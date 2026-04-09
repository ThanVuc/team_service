package adaptermessaginghandler

import (
	"context"
	adaptermessagingconst "team_service/internal/adapter/constant/messaging"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure/share/utils"

	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
)

type AISprintGenerationHandler struct {
	logger        log.LoggerV2
	consumer      eventbus.Consumer
	sprintUseCase usecase.SprintUseCase
}

func NewAISprintGenerationHandler(
	logger log.LoggerV2,
	eventbusConnector *eventbus.RabbitMQConnector,
	sprintUseCase usecase.SprintUseCase,
) *AISprintGenerationHandler {
	consumer := eventbus.NewConsumer(
		eventbusConnector,
		adaptermessagingconst.AISprintGenerationExchangeName,
		eventbus.ExchangeTypeDirect,
		adaptermessagingconst.AISprintGenerationRoutingKey,
		adaptermessagingconst.AISprintGenerationQueueName,
		adaptermessagingconst.AISprintGenerationInstanceNumber,
	)

	return &AISprintGenerationHandler{
		logger:        logger,
		consumer:      consumer,
		sprintUseCase: sprintUseCase,
	}
}

func (h *AISprintGenerationHandler) Handle(ctx context.Context) error {
	return utils.WithSafeMQPanic(h.logger, func() error {
		return h.consumer.Consume(
			ctx,
			h.sprintUseCase.ConsumeAISprintGenerationResult(ctx),
		)
	})
}
