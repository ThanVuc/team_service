package store

import (
	"context"
	irepository "team_service/internal/application/common/interface/repository"
	istore "team_service/internal/application/common/interface/store"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool          *pgxpool.Pool
	repocontainer *repositoryContainer
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:          pool,
		repocontainer: newRepoContainer(pool),
	}
}

func (s *Store) ExecTx(
	ctx context.Context,
	fn func(repo istore.RepositoryContainer) errorbase.AppError,
) (appErr errorbase.AppError) {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	defer func() {
		if appErr != nil {
			tx.Rollback(ctx)
		}
	}()

	repos := newRepoContainer(tx)

	appErr = fn(repos)
	if appErr != nil {
		return
	}

	if err := tx.Commit(ctx); err != nil {
		return errorbase.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (s *Store) GroupRepository() irepository.GroupRepository {
	return s.repocontainer.GroupRepository()
}

func (s *Store) SprintRepository() irepository.SprintRepository {
	return s.repocontainer.SprintRepository()
}

func (s *Store) WorkRepository() irepository.WorkRepository {
	return s.repocontainer.WorkRepository()
}

func (s *Store) UserRepository() irepository.UserRepository {
	return s.repocontainer.UserRepository()
}
