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

	var email *string
	if req.Email != nil {
		email = req.Email
	}

	return &appdto.CreateInviteRequest{
		GroupID: req.GroupId,
		Role:    MapProtoGroupRole(req.Role),
		Email:   email,
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

func ToAcceptInviteRequest(req *team_service.AcceptInviteRequest) *appdto.AcceptInviteRequest {
	if req == nil {
		return nil
	}

	return &appdto.AcceptInviteRequest{
		Code: req.Code,
	}
}

func ToAcceptInviteGrpcResponse(
	resp *appdto.BaseResponse[appdto.AcceptInviteResponse],
) *team_service.AcceptInviteResponse {
	if resp == nil || resp.Data == nil {
		return &team_service.AcceptInviteResponse{
			Location: "",
			Error:    ToProtoError(nil),
		}
	}

	return &team_service.AcceptInviteResponse{
		Location: resp.Data.Location,
		Error:    ToProtoError(resp.Error),
	}
}
