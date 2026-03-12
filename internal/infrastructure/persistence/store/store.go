package store

import (
	"context"
	irepository "team_service/internal/application/common/interface/repository"
	istore "team_service/internal/application/common/interface/store"
	coreerror "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/infrastructure/persistence/db/database"
	"team_service/internal/infrastructure/persistence/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool: pool,
	}
}

func (s *Store) ExecTx(
	ctx context.Context,
	fn func(repo istore.RepositoryContainer) coreerror.AppError,
) coreerror.AppError {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return coreerror.Wrap(err, errdict.ErrInternal)
	}

	defer tx.Rollback(ctx)

	repos := newRepoContainer(tx)

	if err := fn(repos); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return coreerror.Wrap(err, errdict.ErrInternal)
	}

	return nil
}

func (s *Store) GroupRepository() irepository.GroupRepository {
	return repository.NewGroupRepository(database.New(s.pool))
}

func (s *Store) SprintRepository() irepository.SprintRepository {
	return repository.NewSprintRepository(database.New(s.pool))
}

func (s *Store) WorkRepository() irepository.WorkRepository {
	return repository.NewWorkRepository(database.New(s.pool))
}
