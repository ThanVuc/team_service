package repository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/share/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	q *database.Queries
}

func NewUserRepository(
	q *database.Queries,
) *UserRepository {
	return &UserRepository{
		q: q,
	}
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*entity.User, errorbase.AppError) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	u, err := r.q.GetUserByID(ctx, uuid)

	if err != nil {
		return nil, errorbase.Wrap(err, errdict.ErrBadRequest)
	}

	return &entity.User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Status:    enum.UserStatus(u.Status),
		TimeZone:  u.TimeZone,
		CreatedAt: u.CreatedAt.Time,
		AvatarURL: utils.Ptr(u.AvatarUrl.String),
	}, nil
}
