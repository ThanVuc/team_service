package irepository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*entity.User, errorbase.AppError)
}
