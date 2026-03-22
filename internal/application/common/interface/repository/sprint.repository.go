package irepository

import (
	"context"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/domain/entity"
	"team_service/internal/domain/enum"
	"time"
)

type SprintRepository interface {
	DeleteDraftSprint(ctx context.Context, sprintID string) errorbase.AppError
	CancelSprint(ctx context.Context, sprintID string) errorbase.AppError
	DeleteSprint(ctx context.Context, sprintID string) errorbase.AppError
	CreateSprint(ctx context.Context, sprint *entity.Sprint) (*entity.Sprint, errorbase.AppError)
	GetSprintByID(ctx context.Context, sprintID string) (*entity.Sprint, errorbase.AppError)
	UpdateSprint(ctx context.Context, sprintID string, name, goal *string, startDate, endDate *time.Time) (*entity.Sprint, errorbase.AppError)
	UpdateSprintStatus(ctx context.Context, sprintID string, status enum.SprintStatus) (*entity.Sprint, errorbase.AppError)
	IsSprintOverlap(ctx context.Context, groupID string, startDate, endDate time.Time) (bool, errorbase.AppError)
	GetSprintsByGroupID(ctx context.Context, sprintID string) ([]*entity.Sprint, errorbase.AppError)
	GetSimpleSprintsByGroupID(ctx context.Context, groupID string) ([]*appdto.SimpleSprintDTO, errorbase.AppError)
}
