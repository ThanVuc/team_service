package irepository

import (
	"context"
	coreerror "team_service/internal/domain/common/apperror"
	errorbase "team_service/internal/domain/common/apperror"
)

type WorkRepository interface {
	CreateWork() coreerror.AppError
	UnassignWorksByMember(ctx context.Context, groupID string, userID string) errorbase.AppError
}
