package adapter

import (
	grpccontroller "team_service/internal/adapter/gprc"
	"team_service/internal/application"
)

type Dependency struct {
	app *application.Dependency

	GroupController  *grpccontroller.GroupController
	SprintController *grpccontroller.SprintController
	WorkController   *grpccontroller.WorkController
}

func NewDependency(
	applicationDependency *application.Dependency,
) *Dependency {
	return &Dependency{
		app: applicationDependency,
		// ===================================
		// gRPC Controllers
		// ===================================
		GroupController:  grpccontroller.NewGroupController(),
		SprintController: grpccontroller.NewSprintController(),
		WorkController:   grpccontroller.NewWorkController(),

		// ===================================
		// Messaging Handlers
		// ===================================

		// ===================================
		// Job Handlers
		// ===================================
	}
}
