package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/proto/team_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateInviteRequest(req *team_service.CreateInviteRequest) *appdto.CreateInviteRequest {
	if req == nil {
		return nil
	}

	return &appdto.CreateInviteRequest{
		GroupID: req.GroupId,
		Role:    MapProtoGroupRole(req.Role),
		Email:   req.Email,
	}
}

func ToCreateInviteGrpcResponse(
	resp *appdto.BaseResponse[appdto.InviteResponse],
) *team_service.CreateInviteResponse {
	if resp == nil {
		return &team_service.CreateInviteResponse{
			Invite: nil,
			Error:  ToProtoError(nil),
		}
	}

	return &team_service.CreateInviteResponse{
		Invite: ToInviteMessage(resp.Data),
		Error:  ToProtoError(resp.Error),
	}
}

func ToInviteMessage(invite *appdto.InviteResponse) *team_service.InviteMessage {
	if invite == nil {
		return nil
	}

	return &team_service.InviteMessage{
		Code:      invite.Code,
		ExpiresAt: timestamppb.New(invite.ExpiresAt),
		CreatedAt: timestamppb.New(invite.CreatedAt),
	}
}
