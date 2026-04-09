package infrastructure

import (
	"context"
	istore "team_service/internal/application/common/interface/store"
	infracache "team_service/internal/infrastructure/cache"
	cacherepository "team_service/internal/infrastructure/cache/repository"
	grpcinterceptor "team_service/internal/infrastructure/interceptor/grpc"
	"team_service/internal/infrastructure/logging"
	"team_service/internal/infrastructure/messaging"
	"team_service/internal/infrastructure/persistence"
	"team_service/internal/infrastructure/persistence/db"
	"team_service/internal/infrastructure/persistence/store"
	"team_service/internal/infrastructure/r2"
	"team_service/internal/infrastructure/share/settings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanvuc/go-core-lib/cache"
	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
	"github.com/thanvuc/go-core-lib/storage"
	"google.golang.org/grpc"
)

type Dependency struct {
	logger                   log.LoggerV2
	pool                     *pgxpool.Pool
	eventBus                 *eventbus.RabbitMQConnector
	config                   *settings.Configuration
	store                    *store.Store
	recoveryPanicInterceptor grpc.UnaryServerInterceptor
	cacheRepository          *cacherepository.CacheRepository
	cacheClient              *cache.RedisCache
	r2Client                 *storage.R2Client
}

func NewDependency() *Dependency {
	return &Dependency{}
}

func (d *Dependency) Start(ctx context.Context) error {
	err := d.InitInfra(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dependency) InitInfra(ctx context.Context) error {
	// load config
	cfg, err := settings.LoadConfiguration("")
	if err != nil {
		panic(err)
	}
	d.config = cfg

	// logger
	loggerV1 := logging.NewLogger(&cfg.Log)
	d.logger, err = logging.NewLoggerV2()
	if err != nil {
		panic(err)
	}

	// persistence
	d.pool, err = persistence.NewPostgresPool(ctx, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	d.logger.Info("Postgres pool created successfully")
	db.NewMigrations(d.pool, d.logger)

	d.eventBus, err = messaging.NewEventBus(cfg.RabbitMQ, loggerV1)
	if err != nil {
		panic(err)
	}

	// cache
	d.cacheClient, err = infracache.NewRedisClient(cfg.Redis, d.logger)
	if err != nil {
		panic(err)
	}
	d.cacheRepository = cacherepository.NewCacheRepository(d.cacheClient)

	// store
	d.store = store.NewStore(d.pool)

	// r2 client
	d.r2Client, err = r2.NewR2Client(cfg.R2, loggerV1)
	if err != nil {
		panic(err)
	}

	// interceptor
	d.recoveryPanicInterceptor = grpcinterceptor.PanicRecoveryInterceptor(d.logger)

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

	if d.cacheClient != nil {
		err := d.cacheClient.Client.Close()
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

func (d *Dependency) GetStore() istore.Store {
	return d.store
}

func (d *Dependency) GetR2Client() *storage.R2Client {
	return d.r2Client
}

func (d *Dependency) GetRecoveryPanicInterceptor() grpc.UnaryServerInterceptor {
	return d.recoveryPanicInterceptor
}

func (d *Dependency) GetCacheRepository() *cacherepository.CacheRepository {
	return d.cacheRepository
}
