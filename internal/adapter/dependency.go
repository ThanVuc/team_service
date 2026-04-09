package adapter

import (
	"context"
	grpccontroller "team_service/internal/adapter/gprc"
	adaptermessaginghandler "team_service/internal/adapter/messaging"
	"team_service/internal/application"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	// Dependencies
	app   *application.Dependency
	infra *infrastructure.Dependency

	// gRPC Controllers
	GroupController  *grpccontroller.GroupController
	SprintController *grpccontroller.SprintController
	WorkController   *grpccontroller.WorkController
	UserController   *grpccontroller.UserController

	// Messaging Handlers
	AuthHandler               *adaptermessaginghandler.AuthHandler
	AISprintGenerationHandler *adaptermessaginghandler.AISprintGenerationHandler
}

func NewDependency(
	app *application.Dependency,
	infra *infrastructure.Dependency,
) *Dependency {
	return &Dependency{
		app:   app,
		infra: infra,
	}
}

func (d *Dependency) Start(ctx context.Context) error {

	// ===================================
	// gRPC Controllers
	// ===================================

	d.GroupController = grpccontroller.NewGroupController(
		d.app.GroupUseCase,
		d.infra.GetLogger(),
	)

	d.SprintController = grpccontroller.NewSprintController(
		d.app.SprintUseCase,
		d.infra.GetLogger(),
	)
	d.WorkController = grpccontroller.NewWorkController(
		d.app.WorkUseCase,
		d.infra.GetLogger(),
	)

	d.UserController = grpccontroller.NewUserController(
		d.app.UserUseCase,
		d.infra.GetLogger(),
	)
	// ===================================
	// Messaging Handlers
	// ===================================

	d.AuthHandler = adaptermessaginghandler.NewAuthHandler(
		d.infra.GetLogger(),
		d.infra.GetEventBus(),
		d.app.UserUseCase,
	)

	d.AISprintGenerationHandler = adaptermessaginghandler.NewAISprintGenerationHandler(
		d.infra.GetLogger(),
		d.infra.GetEventBus(),
		d.app.SprintUseCase,
	)

	go d.AuthHandler.Handle(ctx)
	go d.AISprintGenerationHandler.Handle(ctx)

	return nil
}

func (d *Dependency) Stop(ctx context.Context) error {
	return nil
}
