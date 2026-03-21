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

type NotificationHelper struct {
	notificationPublisher eventbus.PublisherV2
	logger                log.LoggerV2
}

func NewNotificationHelper(
	eventbusConnector *eventbus.RabbitMQConnector,
	logger log.LoggerV2,
) *NotificationHelper {
	notificationPublisher := eventbus.NewPublisherV2(
		eventbusConnector,
		appconstant.TEAM_NOTIFICATION_EXCHANGE,
		eventbus.ExchangeTypeDirect,
		nil,
		nil,
		false,
		logger,
	)

	return &NotificationHelper{
		notificationPublisher: notificationPublisher,
		logger:                logger,
	}
}

func (h *NotificationHelper) PublishTeamNotificationMessage(
	ctx context.Context,
	message appdto.TeamNotificationMessage,
) errorbase.AppError {
	bytesMessage, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal team notification message", log.WithFields(zap.Error(err)))
		return errorbase.New(errdict.ErrPublishMessage, errorbase.WithDetail("Failed to marshal team notification message."))
	}

	err = h.notificationPublisher.Publish(
		ctx,
		[]string{appconstant.TEAM_NOTIFICATION_ROUTING_KEY},
		bytesMessage,
	)

	if err != nil {
		h.logger.Error("failed to publish team notification message", log.WithFields(zap.Error(err)))
		return errorbase.New(errdict.ErrPublishMessage, errorbase.WithDetail("Failed to publish team notification message to the message broker."))
	}

	return nil
}
