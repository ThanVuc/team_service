package irepository

import (
	"context"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/proto/team_service"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest, userID string) (*database.Group, error)
	CountGroupsByOwner(ctx context.Context, ownerID string) (int64, error)
	GetUserByID(ctx context.Context, userID string) (*database.GetUserByIDRow, error)
	AddGroupMember(ctx context.Context, arg database.CreateGroupMemberParams) error
}
