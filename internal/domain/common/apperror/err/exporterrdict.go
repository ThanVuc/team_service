package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	domainhelper "team_service/internal/domain/common/helper"
)

var (
	ErrSprintExportInvalidInput = errorbase.ErrorInfo{
		Code:   "ts.export.sprint.invalid-input",
		Title:  "Invalid Sprint Export Input",
		Detail: domainhelper.Ptr("The provided sprint export input is invalid."),
	}

	ErrSprintExportGenerateFailed = errorbase.ErrorInfo{
		Code:   "ts.export.sprint.generate-failed",
		Title:  "Sprint Export Generation Failed",
		Detail: domainhelper.Ptr("Failed to generate sprint export file."),
	}
)
