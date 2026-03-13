package appmapper

import (
	appdto "team_service/internal/application/common/dto"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/proto/team_service"
)

func ToGroupResponse(
	group *entity.Group,
	owner *entity.User,
	myRole enum.GroupRole,
	activeSprint *string,
	memberTotal int,
) *appdto.GroupResponse {

	var ownerDTO appdto.OwnerDTO
	if owner != nil {
		ownerDTO = appdto.OwnerDTO{
			ID:     owner.ID,
			Email:  owner.Email,
			Avatar: owner.AvatarURL,
		}
	}

	return &appdto.GroupResponse{
		ID:           group.ID,
		Name:         group.Name,
		Description:  group.Description,
		Owner:        ownerDTO,
		MyRole:       myRole,
		ActiveSprint: activeSprint,
		MemberTotal:  memberTotal,
		AvatarURL:    group.AvatarURL,
		CreatedAt:    group.CreatedAt,
		UpdatedAt:    group.UpdatedAt,
	}
}

func mapSimpleUserMessage(user *entity.User) *team_service.SimpleUserMessage {
	return &team_service.SimpleUserMessage{
		Id:     user.ID,
		Email:  user.Email,
		Avatar: user.AvatarURL,
	}
}

// func MapGroupDetail(group *database.GetGroupByIDRow, owner *entity.User, memberCount int32, role team_service.GroupRole, sprintName string) *team_service.GroupMessage {
// 	ownerMessage := mapSimpleUserMessage(owner)
// 	return &team_service.GroupMessage{
// 		Id:           group.ID.String(),
// 		Name:         group.Name,
// 		Description:  &group.Description.String,
// 		Owner:        ownerMessage,
// 		CreatedAt:    timestamppb.New(group.CreatedAt.Time),
// 		UpdatedAt:    timestamppb.New(group.UpdatedAt.Time),
// 		MemberCount:  memberCount,
// 		MyRole:       role,
// 		ActiveSprint: &sprintName,
// 	}
// }
