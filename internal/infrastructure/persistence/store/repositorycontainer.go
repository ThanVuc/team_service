package store

import (
	irepository "team_service/internal/application/common/interface/repository"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/persistence/repository"

	"github.com/jackc/pgx/v5"
)

type repositoryContainer struct {
	q *database.Queries
}

func newRepoContainer(tx pgx.Tx) *repositoryContainer {
	return &repositoryContainer{
		q: database.New(tx),
	}
}

func (r *repositoryContainer) GroupRepository() irepository.GroupRepository {
	return repository.NewGroupRepository(r.q)
}

func (r *repositoryContainer) SprintRepository() irepository.SprintRepository {
	return repository.NewSprintRepository(r.q)
}

func (r *repositoryContainer) WorkRepository() irepository.WorkRepository {
	return repository.NewWorkRepository(r.q)
}
