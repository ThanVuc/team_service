package utils

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"

	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

func WithSafePanic[TReq any, TResp any](
	ctx context.Context,
	logger log.LoggerV2,
	req TReq,
	f func(context.Context, TReq) (TResp, errorbase.AppError),
) (resp TResp, err error) {

	requestID := GetRequestIDFromOutgoingContext(ctx)

	resp, appErr := f(ctx, req)
	if appErr != nil {

		logger.Error("Usecase returned error",
			log.WithRequestID(requestID),
			log.WithFields(
				zap.Error(appErr),
			),
		)

		err = appErr
	}

	return
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
