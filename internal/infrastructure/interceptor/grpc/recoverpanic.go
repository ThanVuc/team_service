package grpcinterceptor

import (
	"context"
	"team_service/internal/infrastructure/share/utils"

	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecoveryInterceptor(logger log.LoggerV2) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		defer func() {
			if r := recover(); r != nil {

				requestID := utils.GetRequestIDFromOutgoingContext(ctx)

				logger.Error("Recovered from panic",
					log.WithRequestID(requestID),
					log.WithFields(
						zap.Any("panic", r),
						zap.Stack("stacktrace"),
					),
				)

				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
