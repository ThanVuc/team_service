package irepository

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*entity.User, errorbase.AppError)
	GetUserWithPermissionByID(ctx context.Context, userID string, groupId string) (*appdto.UserWithPermission, errorbase.AppError)
	UpsertUser(ctx context.Context, user *entity.User) errorbase.AppError
	GetListMembersByGroupID(ctx context.Context, groupID string) (*appdto.ListMembersResponse, errorbase.AppError)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, errorbase.AppError)
	UpdateUserNotificationSettings(ctx context.Context, userID string, hasEmailNotification bool, hasPushNotification bool) (bool, errorbase.AppError)
}
