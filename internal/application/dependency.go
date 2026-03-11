package application

import (
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	// use cases
	GroupUseCase usecase.GroupUseCase
	UserUseCase  usecase.UserUseCase
}

func NewDependency(infra *infrastructure.Dependency) *Dependency {
	groupUseCase := usecase.NewGroupUseCase(infra.GetStore())
	userUsecase := usecase.NewUserUseCase(infra.GetStore())

	return &Dependency{
		GroupUseCase: groupUseCase,
		UserUseCase:  userUsecase,
	}
}
