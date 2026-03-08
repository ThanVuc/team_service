package store

import (
	"context"
	istore "team_service/internal/application/common/interface/store"

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

func (s *Store) ExecTx(ctx context.Context, fn func(repo istore.RepositoryContainer) error) error {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	repos := newRepoContainer(tx)

	err = fn(repos)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
