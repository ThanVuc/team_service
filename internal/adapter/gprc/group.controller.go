package grpccontroller

import "team_service/proto/team_service"

type GroupController struct {
	team_service.UnimplementedGroupServiceServer
}

func NewGroupController() *GroupController {
	return &GroupController{}
}
