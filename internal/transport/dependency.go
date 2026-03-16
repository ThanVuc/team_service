package transport

import (
	"context"
	"team_service/internal/adapter"
	"team_service/internal/infrastructure"
	transportgrpc "team_service/internal/transport/grpc"
	trasportgrpc "team_service/internal/transport/grpc"
)

type Dependency struct {
	infra   *infrastructure.Dependency
	adapter *adapter.Dependency

	grpcServer *trasportgrpc.TeamServer
}

func NewDependency(
	infra *infrastructure.Dependency,
	adapter *adapter.Dependency,
) *Dependency {
	return &Dependency{
		infra:   infra,
		adapter: adapter,
	}
}

func (d *Dependency) Start(ctx context.Context) error {

	d.grpcServer = transportgrpc.NewTeamServer(
		d.infra.GetLogger(),
		&d.infra.GetConfig().Server,
		d.adapter,
		d.infra.GetRecoveryPanicInterceptor(),
	)

	go d.grpcServer.Start(ctx)

	return nil
}

func (d *Dependency) Stop(ctx context.Context) error {
	if d.grpcServer != nil {
		d.grpcServer.Stop(ctx)
	}
	return nil
}
