package application

import (
	appmapper "team_service/internal/application/common/mapper"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	// use cases
	GroupUseCase usecase.GroupUseCase
	UserUseCase  usecase.UserUseCase
}

func NewDependency(infra *infrastructure.Dependency) *Dependency {
	userUsecase := usecase.NewUserUseCase(infra.GetStore())
	groupUseCase := usecase.NewGroupUseCase(infra.GetStore(), &appmapper.GroupMapper{})

	return &Dependency{
		GroupUseCase: groupUseCase,
		UserUseCase:  userUsecase,
	}
}
