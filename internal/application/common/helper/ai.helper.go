package apphelper

import (
	"context"
	"encoding/json"
	appconstant "team_service/internal/application/common/constant"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"

	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
)

type AIHelper struct {
	aiPublisher eventbus.PublisherV2
	logger      log.LoggerV2
}

func NewAIHelper(
	eventbusConnector *eventbus.RabbitMQConnector,
	logger log.LoggerV2,
) *AIHelper {
	aiPublisher := eventbus.NewPublisherV2(
		eventbusConnector,
		appconstant.AI_TEAM_EXCHANGE,
		eventbus.ExchangeTypeDirect,
		nil,
		nil,
		false,
		logger,
	)

	return &AIHelper{
		aiPublisher: aiPublisher,
		logger:      logger,
	}
}

func (h *AIHelper) PublishSprintGenerationRequest(
	ctx context.Context,
	message appdto.AISprintGenerationRequestedMessage,
) errorbase.AppError {
	bytesMessage, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal ai sprint generation message", log.WithFields(zap.Error(err)))
		return errorbase.New(errdict.ErrPublishMessage, errorbase.WithDetail("Failed to marshal AI sprint generation message."))
	}

	err = h.aiPublisher.Publish(
		ctx,
		[]string{appconstant.AI_TEAM_SPRINT_GENERATION_ROUTING_KEY},
		bytesMessage,
	)
	if err != nil {
		h.logger.Error("failed to publish ai sprint generation message", log.WithFields(zap.Error(err)))
		return errorbase.New(errdict.ErrPublishMessage, errorbase.WithDetail("Failed to publish AI sprint generation message to the message broker."))
	}

	return nil
}
