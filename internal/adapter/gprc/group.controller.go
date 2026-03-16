package grpccontroller

import (
	"context"
	adapermapper "team_service/internal/adapter/mapper"
	"team_service/internal/application/usecase"
	"team_service/internal/infrastructure/share/utils"
	"team_service/proto/common"
	"team_service/proto/team_service"

	"github.com/thanvuc/go-core-lib/log"
)

type GroupController struct {
	team_service.UnimplementedGroupServiceServer
	groupUseCase usecase.GroupUseCase
	logger       log.LoggerV2
}

func NewGroupController(
	groupUseCase usecase.GroupUseCase,
	logger log.LoggerV2,
) *GroupController {
	return &GroupController{
		groupUseCase: groupUseCase,
		logger:       logger,
	}
}

func (c *GroupController) Ping(ctx context.Context, req *common.EmptyRequest) (*common.EmptyResponse, error) {
	return utils.WithSafePanic(ctx, c.logger, req, c.groupUseCase.Ping)
}

func (c *GroupController) CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest) (*team_service.CreateGroupResponse, error) {
	createGroupReq := adapermapper.ToCreateGroupDTO(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createGroupReq, c.groupUseCase.CreateGroup)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateGrouGrpcpResponse(resp), nil
}

func (c *GroupController) GetGroup(ctx context.Context, req *common.IDRequest) (*team_service.GetGroupResponse, error) {
	getGroupReq := adapermapper.ToGetGroupRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getGroupReq, c.groupUseCase.GetGroupRequest)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetGroupGrpcResponse(resp), nil
}
