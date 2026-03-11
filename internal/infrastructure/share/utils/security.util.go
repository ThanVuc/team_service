package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

func WithSafePanic[TReq any, TResp any](
	ctx context.Context,
	logger log.LoggerV2,
	req TReq,
	f func(context.Context, TReq) (TResp, error),
) (resp TResp, err error) {
	requestId := GetRequestIDFromOutgoingContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic",
				log.WithRequestID(requestId),
				log.WithFields(
					zap.Any("panic", r),
					zap.Stack("stacktrace"),
				),
			)
		}
	}()

	return f(ctx, req)
}

func SafeHandler(
	logger log.LoggerV2,
	handler func(d rabbitmq.Delivery) rabbitmq.Action,
) func(d rabbitmq.Delivery) rabbitmq.Action {

	return func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		requestID := GetRequestID(d)

		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic",
					log.WithRequestID(requestID),
					log.WithFields(
						zap.Any("panic", r),
						zap.Stack("stacktrace"),
					),
				)

				action = rabbitmq.NackDiscard
			}
		}()

		action = handler(d)
		return
	}
}

func RetryConsumer(
	ctx context.Context,
	logger log.LoggerV2,
	retryDelay time.Duration,
	name string,
	run func(ctx context.Context) error,
) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						fmt.Sprintf("%s consumer panic recovered", name),
						log.WithFields(zap.Any("panic", r)),
					)
				}
			}()

			if err := run(ctx); err != nil {
				logger.Error(
					fmt.Sprintf("%s consumer stopped", name),
					log.WithFields(zap.Error(err)),
				)
			}
		}()

		select {
		case <-ctx.Done():
			logger.Info(fmt.Sprintf("%s consumer stopped by context", name))
			return
		case <-time.After(retryDelay):
		}
	}
}
