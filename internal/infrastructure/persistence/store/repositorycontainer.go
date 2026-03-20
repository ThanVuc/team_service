package store

import (
	irepository "team_service/internal/application/common/interface/repository"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/persistence/repository"
)

type repositoryContainer struct {
	q *database.Queries
}

func newRepoContainer(db database.DBTX) *repositoryContainer {
	return &repositoryContainer{
		q: database.New(db),
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

func (r *repositoryContainer) UserRepository() irepository.UserRepository {
	return repository.NewUserRepository(r.q)
}

func (r *repositoryContainer) InviteRepository() irepository.InviteRepository {
	return repository.NewInviteRepository(r.q)
}
