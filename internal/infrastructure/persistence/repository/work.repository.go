package repository

import "team_service/internal/infrastructure/persistence/db/database"

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

func (r *WorkRepository) CreateWork() error {
	// Implement the logic to create a work item in the database
	return nil
}
