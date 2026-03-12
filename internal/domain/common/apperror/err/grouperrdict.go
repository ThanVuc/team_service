package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	"team_service/internal/infrastructure/share/utils"
)

var (
	ErrInvalidUUID = errorbase.ErrorInfo{
		Code:   "ts.validation.invalid-uuid",
		Title:  "Invalid UUID",
		Detail: utils.Ptr("The provided UUID is not valid. Please ensure it is in the correct format."),
	}
)
