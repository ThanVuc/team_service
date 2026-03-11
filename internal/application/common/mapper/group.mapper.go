package appmapper

import (
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/proto/team_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GroupMapper struct{}

func (m *GroupMapper) MapGroupMessage(group *database.Group, user *database.GetUserByIDRow) *team_service.GroupMessage {
	userMessage := m.mapSimpleUserMessage(user)
	return &team_service.GroupMessage{
		Id:          group.ID.String(),
		Name:        group.Name,
		Description: group.Description.String,
		Owner:       userMessage,
		CreatedAt:   timestamppb.New(group.CreatedAt.Time),
		UpdatedAt:   timestamppb.New(group.UpdatedAt.Time),
	}
}

func (m *GroupMapper) mapSimpleUserMessage(db *database.GetUserByIDRow) *team_service.SimpleUserMessage {
	return &team_service.SimpleUserMessage{
		Id:     db.ID.String(),
		Email:  db.Email,
		Avatar: db.AvatarUrl.String,
	}
}
