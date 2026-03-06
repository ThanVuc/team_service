package infrastructure

import (
	"context"
	"team_service/internal/infrastructure/logging"
	"team_service/internal/infrastructure/messaging"
	"team_service/internal/infrastructure/persistence"
	"team_service/internal/infrastructure/share/settings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
)

type Dependency struct {
	logger   log.LoggerV2
	pool     *pgxpool.Pool
	eventBus *eventbus.RabbitMQConnector
	config   *settings.Configuration
}

func NewDependency() *Dependency {
	return &Dependency{}
}

func (d *Dependency) Start(ctx context.Context) error {
	cfg, err := settings.LoadConfiguration("")
	if err != nil {
		panic(err)
	}
	d.config = cfg

	loggerV1 := logging.NewLogger(&cfg.Log)
	d.logger, err = logging.NewLoggerV2()
	if err != nil {
		panic(err)
	}

	d.pool, err = persistence.NewPostgresPool(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	d.logger.Info("Postgres pool created successfully")

	d.eventBus, err = messaging.NewEventBus(cfg.RabbitMQ, loggerV1)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d *Dependency) Stop(ctx context.Context) error {
	if d.pool != nil {
		d.pool.Close()
	}

	if d.eventBus != nil {
		err := d.eventBus.Shutdown(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dependency) GetLogger() log.LoggerV2 {
	return d.logger
}

func (d *Dependency) GetDBPool() *pgxpool.Pool {
	return d.pool
}

func (d *Dependency) GetEventBus() *eventbus.RabbitMQConnector {
	return d.eventBus
}

func (d *Dependency) GetConfig() *settings.Configuration {
	return d.config
}
