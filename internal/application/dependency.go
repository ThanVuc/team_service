package application

import (
	"context"

	apphelper "team_service/internal/application/common/helper"
	appvalidation "team_service/internal/application/common/validation"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	infra *infrastructure.Dependency

	// use cases
	GroupUseCase usecase.GroupUseCase
	UserUseCase  usecase.UserUseCase
}

func NewDependency(infra *infrastructure.Dependency) *Dependency {
	return &Dependency{
		infra: infra,
	}
}

func (d *Dependency) Start(ctx context.Context) error {
	return d.InitUseCases(ctx)
}

func (d *Dependency) InitUseCases(ctx context.Context) error {
	store := d.infra.GetStore()

	// ============ validators ============
	groupValidator := appvalidation.NewGroupValidator(
		store.GroupRepository(),
		store.UserRepository(),
	)

	// ============ helper ============
	authHelper := apphelper.NewAuthHelper(
		store.UserRepository(),
		d.infra.GetCacheRepository(),
		d.infra.GetLogger(),
	)

	// ============ use cases ============
	d.UserUseCase = usecase.NewUserUseCase(store, d.infra.GetLogger())
	d.GroupUseCase = usecase.NewGroupUseCase(
		store,
		groupValidator,
		authHelper,
	)

	return nil
}

func (d *Dependency) Stop(ctx context.Context) error {
	// usually nothing to stop in application layer
	return nil
}
