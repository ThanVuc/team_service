package irepository

import (
	"context"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
	"team_service/internal/infrastructure/persistence/db/database"
	"time"
)

type SprintRepository interface {
	DeleteDraftSprint(ctx context.Context, sprintID string) errorbase.AppError
	CancelSprint(ctx context.Context, sprintID string) errorbase.AppError
	GetSprintsByGroupID(ctx context.Context, sprintID string) ([]database.Sprint, errorbase.AppError)
	CreateSprint(ctx context.Context, sprint *entity.Sprint) (*entity.Sprint, errorbase.AppError)
	IsSprintOverlap(ctx context.Context, groupID string, startDate, endDate time.Time) (bool, errorbase.AppError)
}
