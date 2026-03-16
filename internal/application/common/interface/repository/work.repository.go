package irepository

import coreerror "team_service/internal/domain/common/apperror"

type WorkRepository interface {
	CreateWork() coreerror.AppError
}
