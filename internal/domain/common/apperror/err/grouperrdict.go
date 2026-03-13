package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	domainhelper "team_service/internal/domain/common/helper"
)

var (
	ErrInvalidUUID = errorbase.ErrorInfo{
		Code:   "ts.validation.invalid-uuid",
		Title:  "Invalid UUID",
		Detail: domainhelper.Ptr("The provided UUID is not valid. Please ensure it is in the correct format."),
	}
)
