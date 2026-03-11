package adaptermessaginghandler

import (
	"context"
	adaptermessagingconst "team_service/internal/adapter/constant/messaging"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure/share/utils"

	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
)

type AuthHandler struct {
	logger      log.LoggerV2
	consumer    eventbus.Consumer
	userUseCase usecase.UserUseCase
}

func NewAuthHandler(
	logger log.LoggerV2,
	eventbusConnector *eventbus.RabbitMQConnector,
	userUseCase usecase.UserUseCase,
) *AuthHandler {
	consumer := eventbus.NewConsumer(
		eventbusConnector,
		adaptermessagingconst.AuthExchangeName,
		eventbus.ExchangeTypeTopic,
		adaptermessagingconst.AuthRoutingKey,
		adaptermessagingconst.AuthQueueName,
		adaptermessagingconst.InstanceNumber,
	)

	return &AuthHandler{
		logger:      logger,
		consumer:    consumer,
		userUseCase: userUseCase,
	}
}

func (h *AuthHandler) Handle(ctx context.Context) {
	utils.RetryConsumer(
		ctx,
		h.logger,
		adaptermessagingconst.RetryInterval,
		adaptermessagingconst.HandlerName,
		func(ctx context.Context) error {
			return h.consumer.Consume(
				ctx,
				utils.SafeHandler(h.logger, h.userUseCase.SyncUserData(ctx)),
			)
		},
	)
}
