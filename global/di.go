package global

import (
	"context"
	"team_service/internal/adapter"
	"team_service/internal/application"
	"team_service/internal/infrastructure"
	"team_service/internal/transport"
)

type GlobalDependency struct {
	infra     *infrastructure.Dependency
	app       *application.Dependency
	adapter   *adapter.Dependency
	transport *transport.Dependency

	lifecycles []Lifecycle
}

func NewGlobalDependency() *GlobalDependency {
	infra := infrastructure.NewDependency()
	app := application.NewDependency(infra)
	adapter := adapter.NewDependency(app, infra)
	transport := transport.NewDependency(infra, adapter)

	g := &GlobalDependency{
		infra:     infra,
		transport: transport,
		app:       app,
		adapter:   adapter,
	}

	g.register(
		infra,
		app,
		transport,
		adapter,
	)

	return g
}

func (g *GlobalDependency) register(components ...Lifecycle) {
	g.lifecycles = append(g.lifecycles, components...)
}

func (g *GlobalDependency) Start(ctx context.Context) error {
	for _, c := range g.lifecycles {
		if err := c.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (g *GlobalDependency) Stop(ctx context.Context) error {
	for i := len(g.lifecycles) - 1; i >= 0; i-- {
		if err := g.lifecycles[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
