package istore

import (
	"context"
	irepository "team_service/internal/application/common/interface/repository"
	errorbase "team_service/internal/domain/common/apperror"
)

type RepositoryContainer interface {
	GroupRepository() irepository.GroupRepository
	SprintRepository() irepository.SprintRepository
	WorkRepository() irepository.WorkRepository
}

type Store interface {
	ExecTx(ctx context.Context, fn func(repo RepositoryContainer) errorbase.AppError) errorbase.AppError
}
