package repository

import (
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/infrastructure/persistence/db/database"
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
