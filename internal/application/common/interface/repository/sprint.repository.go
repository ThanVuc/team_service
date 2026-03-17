package irepository

import (
	"context"
	coreerror "team_service/internal/domain/common/apperror"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/infrastructure/persistence/db/database"
)

type SprintRepository interface {
	CreateSprint() coreerror.AppError
	DeleteDraftSprint(ctx context.Context, sprintID string) errorbase.AppError
	CancelSprint(ctx context.Context, sprintID string) errorbase.AppError
	GetSprintsByGroupID(ctx context.Context, sprintID string) ([]database.Sprint, errorbase.AppError)
}
