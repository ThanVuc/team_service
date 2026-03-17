package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/enum"
	"team_service/proto/common"
	"team_service/proto/team_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateGroupDTO(req *team_service.CreateGroupRequest) *appdto.CreateGroupRequest {
	return &appdto.CreateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
	}
}

func ToCreateGrouGrpcpResponse(
	group *appdto.BaseResponse[appdto.GroupResponse],
) *team_service.CreateGroupResponse {
	if group == nil {
		return &team_service.CreateGroupResponse{
			Group: nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.CreateGroupResponse{
		Group: ToGroupMessage(group.Data),
		Error: ToProtoError(group.Error),
	}
}

func MapEnumGroupRoleToGroupRole(role enum.GroupRole) team_service.GroupRole {
	switch role {
	case enum.GroupRoleOwner:
		return team_service.GroupRole_GROUP_ROLE_OWNER
	case enum.GroupRoleManager:
		return team_service.GroupRole_GROUP_ROLE_MANAGER
	case enum.GroupRoleMember:
		return team_service.GroupRole_GROUP_ROLE_MEMBER
	default:
		return team_service.GroupRole_GROUP_ROLE_VIEWER
	}
}
func ToGroupMessage(group *appdto.GroupResponse) *team_service.GroupMessage {
	if group == nil {
		return nil
	}

	var activeSprintId string
	if group.ActiveSprint != nil {
		activeSprintId = *group.ActiveSprint
	}

	var avatarUrl string
	if group.AvatarURL != nil {
		avatarUrl = *group.AvatarURL
	}

	return &team_service.GroupMessage{
		Id:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Owner: &team_service.SimpleUserMessage{
			Id:     group.Owner.ID,
			Email:  group.Owner.Email,
			Avatar: group.Owner.Avatar,
		},
		Avatar:       avatarUrl,
		ActiveSprint: &activeSprintId,
		MyRole:       MapEnumGroupRoleToGroupRole(group.MyRole),
		MemberCount:  int32(group.MemberTotal),
		CreatedAt:    timestamppb.New(group.CreatedAt),
		UpdatedAt:    timestamppb.New(group.UpdatedAt),
	}
}

func ToGetGroupRequest(req *common.IDRequest) *appdto.GetGroupRequest {
	return &appdto.GetGroupRequest{
		GroupID: req.Id,
	}
}

func ToGetGroupGrpcResponse(
	group *appdto.BaseResponse[appdto.GroupResponse],
) *team_service.GetGroupResponse {
	if group == nil {
		return &team_service.GetGroupResponse{
			Group: nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.GetGroupResponse{
		Group: ToGroupMessage(group.Data),
		Error: ToProtoError(group.Error),
	}
}

func ToUpdateGroupRequest(req *team_service.UpdateGroupRequest) *appdto.UpdateGroupRequest {
	return &appdto.UpdateGroupRequest{
		GroupID:     req.Id,
		Name:        req.Name,
		Description: req.Description,
	}

}

func ToUpdateGroupGrpcResponse(
	group *appdto.BaseResponse[appdto.GroupResponse],
) *team_service.UpdateGroupResponse {
	if group == nil {
		return &team_service.UpdateGroupResponse{
			Group: nil,
			Error: ToProtoError(nil),
		}
	}

	return &team_service.UpdateGroupResponse{
		Group: ToGroupMessage(group.Data),
		Error: ToProtoError(group.Error),
	}
}

func ToDeleteGroupRequest(req *common.IDRequest) *appdto.DeleteGroupRequest {
	return &appdto.DeleteGroupRequest{
		GroupID: req.Id,
	}
}

func ToDeleteGroupGrpcResponse(
	resp *appdto.BaseResponse[appdto.DeleteGroupResponse],
) *team_service.DeleteGroupResponse {

	if resp == nil || resp.Data == nil {
		return &team_service.DeleteGroupResponse{
			Success: false,
			Error:   ToProtoError(resp.Error),
		}
	}

	return &team_service.DeleteGroupResponse{
		Success: resp.Data.Success,
		Error:   ToProtoError(resp.Error),
	}
}
