package irepository

import coreerror "team_service/internal/domain/common/apperror"

type SprintRepository interface {
	CreateSprint() coreerror.AppError
}
