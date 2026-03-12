package repository

import (
	"team_service/internal/infrastructure/persistence/db/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryContainer struct {
	GroupRepository  *GroupRepository
	SprintRepository *SprintRepository
	WorkRepository   *WorkRepository
}

func NewRepositoryContainer(
	pool *pgxpool.Pool,
) *RepositoryContainer {
	return &RepositoryContainer{
		GroupRepository:  NewGroupRepository(database.New(pool)),
		SprintRepository: NewSprintRepository(database.New(pool)),
		WorkRepository:   NewWorkRepository(database.New(pool)),
	}
}

func (c *RepositoryContainer) GetGroupRepository() *GroupRepository {
	return c.GroupRepository
}

func (c *RepositoryContainer) GetSprintRepository() *SprintRepository {
	return c.SprintRepository
}

func (c *RepositoryContainer) GetWorkRepository() *WorkRepository {
	return c.WorkRepository
}
