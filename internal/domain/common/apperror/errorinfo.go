package coreerror


type AppError interface {
	Error() ErrorInfo
}

type ErrorInfo struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

var (
	ErrBadRequest = ErrorInfo{
		Code:  "ts.validation.bad-request",
		Title: "Request invalid",
	}

	ErrUnauthorized = ErrorInfo{
		Code:  "ts.auth.unauthorized",
		Title: "Unauthorized",
	}

	ErrForbidden = ErrorInfo{
		Code:  "ts.auth.forbidden",
		Title: "Forbidden",
	}

	ErrNotFound = ErrorInfo{
		Code:  "ts.resource.not-found",
		Title: "Not Found",
	}

	ErrConflict = ErrorInfo{
		Code:  "ts.resource.conflict",
		Title: "Conflict",
	}

	ErrUnprocessable = ErrorInfo{
		Code:  "ts.validation.unprocessable",
		Title: "Unprocessable",
	}

	ErrInternal = ErrorInfo{
		Code:  "ts.internal.error",
		Title: "Internal Server Error",
	}
)
