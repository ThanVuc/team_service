package messaging

import (
	"fmt"
	"team_service/internal/infrastructure/share/settings"

	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
)

func NewEventBus(cfg settings.RabbitMQ, logger log.Logger) (*eventbus.RabbitMQConnector, error) {
	uri := fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	connector, err := eventbus.NewConnector(uri, logger)
	if err != nil {
		return nil, err
	}

	logger.Info("Event bus created successfully", "")

	return connector, nil
}
