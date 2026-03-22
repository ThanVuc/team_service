package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
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
		MyRole:       MapGroupRole(group.MyRole),
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

func ToListMembersRequest(req *team_service.ListMembersRequest) *appdto.ListMembersRequest {
	return &appdto.ListMembersRequest{
		GroupID: req.GroupId,
	}
}

func ToGetSimpleUserByGroupIDRequest(req *common.IDRequest) *appdto.ListMembersRequest {
	if req == nil {
		return nil
	}

	return &appdto.ListMembersRequest{
		GroupID: req.Id,
	}
}

func ToListMembersGrpcResponse(
	resp *appdto.BaseResponse[appdto.ListMembersResponse],
) *team_service.ListMembersResponse {
	if resp == nil || resp.Data == nil {
		return &team_service.ListMembersResponse{
			Members: nil,
			Error:   ToProtoError(resp.Error),
		}
	}

	members := make([]*team_service.MemberMessage, len(resp.Data.Members))
	for i, member := range resp.Data.Members {
		avatar := ""
		if member.Avatar != nil {
			avatar = *member.Avatar
		}
		members[i] = &team_service.MemberMessage{
			Id:       member.ID,
			Email:    member.Email,
			Avatar:   avatar,
			Role:     MapGroupRole(member.Role),
			JoinedAt: timestamppb.New(member.JoinedAt),
		}
	}

	return &team_service.ListMembersResponse{
		Members: members,
		Error:   ToProtoError(resp.Error),
	}

}

func ToUpdateMemberRoleRequest(req *team_service.UpdateMemberRoleRequest) *appdto.UpdateMemberRoleRequest {
	return &appdto.UpdateMemberRoleRequest{
		GroupID:  req.GroupId,
		MemberId: req.MemberId,
		Role:     MapProtoGroupRole(req.NewRole),
	}
}

func ToUpdateMemberRoleGrpcResponse(
	resp *appdto.BaseResponse[appdto.MemberResponse],
) *team_service.UpdateMemberRoleResponse {
	if resp == nil || resp.Data == nil {
		return &team_service.UpdateMemberRoleResponse{
			Member: nil,
			Error:  ToProtoError(resp.Error),
		}
	}
	avatar := ""
	if resp.Data.Avatar != nil {
		avatar = *resp.Data.Avatar
	}

	return &team_service.UpdateMemberRoleResponse{
		Member: &team_service.MemberMessage{
			Id:       resp.Data.ID,
			Email:    resp.Data.Email,
			Avatar:   avatar,
			Role:     MapGroupRole(resp.Data.Role),
			JoinedAt: timestamppb.New(resp.Data.JoinedAt),
		},
		Error: ToProtoError(resp.Error),
	}
}

func ToRemoveMemberRequest(req *team_service.RemoveMemberRequest) *appdto.RemoveMemberRequest {
	return &appdto.RemoveMemberRequest{
		GroupID:  req.GroupId,
		MemberId: req.MemberId,
	}
}

func ToRemoveMemberGrpcResponse(
	resp *appdto.BaseResponse[appdto.RemoveMemberResponse],
) *team_service.RemoveMemberResponse {
	if resp == nil || resp.Data == nil {
		return &team_service.RemoveMemberResponse{
			Success: false,
			Error:   ToProtoError(resp.Error),
		}
	}

	return &team_service.RemoveMemberResponse{
		Success: resp.Data.Success,
		Error:   ToProtoError(resp.Error),
	}
}

func ToGetSimpleUserByGroupIDGrpcResponse(
	resp *appdto.BaseResponse[[]appdto.SimpleUserResponse],
) *team_service.GetSimpleUserByGroupIDResponse {
	if resp == nil {
		return &team_service.GetSimpleUserByGroupIDResponse{
			Users: nil,
			Error: ToProtoError(nil),
		}
	}

	users := make([]*team_service.SimpleUserMessage, 0)
	if resp.Data != nil {
		users = make([]*team_service.SimpleUserMessage, 0, len(*resp.Data))
		for _, user := range *resp.Data {
			avatar := ""
			if user.AvatarURL != nil {
				avatar = *user.AvatarURL
			}

			users = append(users, &team_service.SimpleUserMessage{
				Id:     user.ID,
				Email:  user.Email,
				Avatar: &avatar,
			})
		}
	}

	return &team_service.GetSimpleUserByGroupIDResponse{
		Users: users,
		Error: ToProtoError(resp.Error),
	}
}
