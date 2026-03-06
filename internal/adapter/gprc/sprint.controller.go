package grpccontroller

import "team_service/proto/team_service"

type SprintController struct {
	team_service.UnimplementedSprintServiceServer
}

func NewSprintController() *SprintController {
	return &SprintController{}
}
