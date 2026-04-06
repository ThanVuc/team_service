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
	userUseCase  usecase.UserUseCase
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

	resp, err := utils.WithSafePanic(ctx, c.logger, getGroupReq, c.groupUseCase.GetGroup)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetGroupGrpcResponse(resp), nil
}

func (c *GroupController) ListGroups(ctx context.Context, req *common.IDRequest) (*team_service.ListGroupsResponse, error) {
	listGroupsReq := adapermapper.ToListGroupsRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, listGroupsReq, c.groupUseCase.ListGroups)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToListGroupsGrpcResponse(resp), nil
}

func (c *GroupController) UpdateGroup(ctx context.Context, req *team_service.UpdateGroupRequest) (*team_service.UpdateGroupResponse, error) {
	updateGroupReq := adapermapper.ToUpdateGroupRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateGroupReq, c.groupUseCase.UpdateGroup)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateGroupGrpcResponse(resp), nil
}

func (c *GroupController) DeleteGroup(ctx context.Context, req *common.IDRequest) (*team_service.DeleteGroupResponse, error) {
	deleteGroupReq := adapermapper.ToDeleteGroupRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, deleteGroupReq, c.groupUseCase.DeleteGroup)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToDeleteGroupGrpcResponse(resp), nil
}

func (c *GroupController) ListMembers(ctx context.Context, req *team_service.ListMembersRequest) (*team_service.ListMembersResponse, error) {
	listMembersReq := adapermapper.ToListMembersRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, listMembersReq, c.groupUseCase.GetListGroupMembers)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToListMembersGrpcResponse(resp), nil
}

func (c *GroupController) GetSimpleUserByGroupID(ctx context.Context, req *common.IDRequest) (*team_service.GetSimpleUserByGroupIDResponse, error) {
	getSimpleUsersReq := adapermapper.ToGetSimpleUserByGroupIDRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, getSimpleUsersReq, c.groupUseCase.GetSimpleUserByGroupID)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGetSimpleUserByGroupIDGrpcResponse(resp), nil
}

func (c *GroupController) UpdateMemberRole(ctx context.Context, req *team_service.UpdateMemberRoleRequest) (*team_service.UpdateMemberRoleResponse, error) {
	updateMemberRoleReq := adapermapper.ToUpdateMemberRoleRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, updateMemberRoleReq, c.groupUseCase.UpdateMemberRole)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToUpdateMemberRoleGrpcResponse(resp), nil
}

func (c *GroupController) RemoveMember(ctx context.Context, req *team_service.RemoveMemberRequest) (*team_service.RemoveMemberResponse, error) {
	removeMemberReq := adapermapper.ToRemoveMemberRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, removeMemberReq, c.groupUseCase.RemoveMember)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToRemoveMemberGrpcResponse(resp), nil
}

func (c *GroupController) CreateInvite(ctx context.Context, req *team_service.CreateInviteRequest) (*team_service.CreateInviteResponse, error) {
	createInviteReq := adapermapper.ToCreateInviteRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, createInviteReq, c.groupUseCase.CreateInvite)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToCreateInviteGrpcResponse(resp), nil
}

func (c *GroupController) AcceptInvite(ctx context.Context, req *team_service.AcceptInviteRequest) (*team_service.AcceptInviteResponse, error) {
	acceptInviteReq := adapermapper.ToAcceptInviteRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, acceptInviteReq, c.groupUseCase.AcceptInvite)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToAcceptInviteGrpcResponse(resp), nil
}

func (c *GroupController) GeneratePresignedURLs(ctx context.Context, req *team_service.GeneratePresignedURLsRequest) (*team_service.GeneratePresignedURLsResponse, error) {
	generatePresignedURLsReq := adapermapper.ToGeneratePresignedURLsRequest(req)

	resp, err := utils.WithSafePanic(ctx, c.logger, generatePresignedURLsReq, c.groupUseCase.GeneratePresignedURLs)
	if err != nil {
		return nil, err
	}

	return adapermapper.ToGeneratePresignedURLsGrpcResponse(resp), nil
}
