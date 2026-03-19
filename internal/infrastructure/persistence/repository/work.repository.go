package repository

import (
	"context"
	"fmt"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/infrastructure/persistence/db/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type WorkRepository struct {
	q *database.Queries
}

func NewWorkRepository(
	q *database.Queries,
) *WorkRepository {
	return &WorkRepository{
		q: q,
	}
}

func (r *WorkRepository) CreateWork() errorbase.AppError {
	// Implement the logic to create a work item in the database
	return nil
}

func (r *WorkRepository) UnassignWorksByMember(ctx context.Context, groupID string, userID string) errorbase.AppError {
	var groupUUID pgtype.UUID
	if err := groupUUID.Scan(groupID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse user id"),
		)
	}

	err := r.q.UnassignWorksByMember(ctx, database.UnassignWorksByMemberParams{
		GroupID:    groupUUID,
		AssigneeID: userUUID,
	})

	if err != nil {
		return errorbase.Wrap(
			err,
			errdict.ErrInternal,
			errorbase.WithDetail(fmt.Sprintf("failed to remove member user=%s from group=%s", userID, groupID)),
		)
	}

	return nil
}
