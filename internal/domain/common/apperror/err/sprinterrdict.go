package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	domainhelper "team_service/internal/domain/common/helper"
)

var (
	StatusTransitionInvalid = errorbase.ErrorInfo{
		Code:   "ts.validation.sprint.invalid-status-transition",
		Title:  "Invalid Sprint Status Transition",
		Detail: domainhelper.Ptr("The requested status transition for the sprint is not allowed. Please ensure the transition follows the defined workflow."),
	}
)
