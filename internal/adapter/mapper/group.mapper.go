package adapermapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/proto/team_service"
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

	return &team_service.GroupMessage{
		Id:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Owner: &team_service.SimpleUserMessage{
			Id:     group.Owner.ID,
			Email:  group.Owner.Email,
			Avatar: group.Owner.Avatar,
		},
		Avatar: *group.AvatarURL,
	}
}
