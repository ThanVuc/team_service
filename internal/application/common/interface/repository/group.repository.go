package irepository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *entity.Group, userID string) (*entity.Group, errorbase.AppError)
	CountGroupsByOwner(ctx context.Context, ownerID string) (int64, errorbase.AppError)
	AddGroupMember(ctx context.Context, member *entity.GroupMember) errorbase.AppError
	GetGroupByID(ctx context.Context, groupID string) (*entity.Group, int32, string, errorbase.AppError)
	GetRoleByUserIDAndGroupID(ctx context.Context, userID string, groupID string) (string, errorbase.AppError)
	CheckGroupExists(ctx context.Context, groupID string) (bool, errorbase.AppError)
	UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, errorbase.AppError)
	DeleteGroup(ctx context.Context, groupID string) errorbase.AppError
}
