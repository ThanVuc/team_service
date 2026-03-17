package repository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/infrastructure/persistence/db/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type SprintRepository struct {
	q *database.Queries
}

func NewSprintRepository(
	q *database.Queries,
) *SprintRepository {
	return &SprintRepository{
		q: q,
	}
}

func (r *SprintRepository) CreateSprint() errorbase.AppError {
	// Implement the logic to create a sprint in the database
	return nil
}

func (r *SprintRepository) CancelSprint(ctx context.Context, sprintID string) errorbase.AppError {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	err := r.q.CancelActiveSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (r *SprintRepository) DeleteDraftSprint(ctx context.Context, sprintID string) errorbase.AppError {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	err := r.q.DeleteDraftSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (r *SprintRepository) GetSprintsByGroupID(
	ctx context.Context,
	sprintID string,
) ([]database.Sprint, errorbase.AppError) {
	var sprintUUID pgtype.UUID
	if err := sprintUUID.Scan(sprintID); err != nil {
		return nil, errorbase.New(
			errdict.ErrInternal,
			errorbase.WithDetail("failed to parse group id"),
		)
	}

	sprints, err := r.q.GetSprintsByGroupID(ctx, sprintUUID)
	if err != nil {
		return nil, errorbase.Wrap(err, errdict.ErrInternal)
	}

	return sprints, nil
}
