package adapter

import (
	grpccontroller "team_service/internal/adapter/gprc"
	"team_service/internal/application"
	"team_service/internal/infrastructure"
)

type Dependency struct {
	// gRPC Controllers
	GroupController  *grpccontroller.GroupController
	SprintController *grpccontroller.SprintController
	WorkController   *grpccontroller.WorkController
}

func NewDependency(
	applicationDependency *application.Dependency,
	infrastructureDependency *infrastructure.Dependency,
) *Dependency {
	// ===================================
	// gRPC Controllers
	// ===================================
	groupController := grpccontroller.NewGroupController(
		applicationDependency.GroupUseCase,
		infrastructureDependency.GetLogger(),
	)
	sprintController := grpccontroller.NewSprintController()
	workController := grpccontroller.NewWorkController()

	// ===================================
	// Messaging Handlers
	// ===================================

	// ===================================
	// Job Handlers
	// ===================================

	return &Dependency{
		// ===================================
		// gRPC Controllers
		// ===================================
		GroupController:  groupController,
		SprintController: sprintController,
		WorkController:   workController,

		// ===================================
		// Messaging Handlers
		// ===================================

		// ===================================
		// Job Handlers
		// ===================================
	}
}
