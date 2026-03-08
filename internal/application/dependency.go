package application

import (
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	// use cases
	GroupUseCase usecase.GroupUseCase
}

func NewDependency(infra *infrastructure.Dependency) *Dependency {
	groupUseCase := usecase.NewGroupUseCase(infra.GetStore())

	return &Dependency{
		GroupUseCase: groupUseCase,
	}
}
