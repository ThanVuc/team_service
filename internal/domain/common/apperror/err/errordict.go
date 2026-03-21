package errdict

import (
	errorbase "team_service/internal/domain/common/apperror"
	domainhelper "team_service/internal/domain/common/helper"
)

// Common errors
var (
	ErrBadRequest = errorbase.ErrorInfo{
		Code:   "ts.validation.bad-request",
		Title:  "Request Invalid",
		Detail: domainhelper.Ptr("The request parameters or body are invalid."),
	}

	ErrUnauthorized = errorbase.ErrorInfo{
		Code:   "ts.auth.unauthorized",
		Title:  "Unauthorized",
		Detail: domainhelper.Ptr("Authentication is required and has failed or has not yet been provided."),
	}

	ErrForbidden = errorbase.ErrorInfo{
		Code:   "ts.auth.forbidden",
		Title:  "Forbidden",
		Detail: domainhelper.Ptr("You do not have permission to access this resource."),
	}

	ErrNotFound = errorbase.ErrorInfo{
		Code:   "ts.resource.not-found",
		Title:  "Not Found",
		Detail: domainhelper.Ptr("The requested resource could not be found."),
	}

	ErrConflict = errorbase.ErrorInfo{
		Code:   "ts.resource.conflict",
		Title:  "Conflict",
		Detail: domainhelper.Ptr("The request could not be completed due to a conflict with the current state of the resource."),
	}

	ErrUnprocessable = errorbase.ErrorInfo{
		Code:   "ts.validation.unprocessable",
		Title:  "Unprocessable Entity",
		Detail: domainhelper.Ptr("The request was well-formed but contains semantic errors."),
	}

	ErrInternal = errorbase.ErrorInfo{
		Code:   "ts.internal.error",
		Title:  "Internal Server Error",
		Detail: domainhelper.Ptr("An unexpected error occurred on the server."),
	}

	ErrPublishMessage = errorbase.ErrorInfo{
		Code:   "ts.messaging.publish-error",
		Title:  "Message Publish Error",
		Detail: domainhelper.Ptr("Failed to publish message to the message broker."),
	}
)
