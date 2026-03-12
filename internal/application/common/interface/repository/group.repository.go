package irepository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/proto/team_service"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, req *team_service.CreateGroupRequest, userID string) (*database.Group, errorbase.AppError)
	CountGroupsByOwner(ctx context.Context, ownerID string) (int64, errorbase.AppError)
	GetUserByID(ctx context.Context, userID string) (*database.GetUserByIDRow, errorbase.AppError)
	AddGroupMember(ctx context.Context, arg database.CreateGroupMemberParams) errorbase.AppError
	GetGroupByID(ctx context.Context, user, groupID string) (*database.GetGroupByIDRow, int32, string, string, errorbase.AppError)
}
