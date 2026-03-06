package grpccontroller

import "team_service/proto/team_service"

type WorkController struct {
	team_service.UnimplementedWorkServiceServer
}

func NewWorkController() *WorkController {
	return &WorkController{}
}
