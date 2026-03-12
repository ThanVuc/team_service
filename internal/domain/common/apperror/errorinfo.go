package errorbase

type appError struct {
	info  ErrorInfo
	cause error
}

type ErrorInfo struct {
	Code   string  `json:"code"`
	Title  string  `json:"title"`
	Detail *string `json:"detail,omitempty"`
}

func (e *appError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.info.Title
}

func (e *appError) ErrorInfo() ErrorInfo {
	return e.info
}

func (e *appError) Unwrap() error {
	return e.cause
}

func New(info ErrorInfo) AppError {
	return &appError{
		info: info,
	}
}

func Wrap(err error, info ErrorInfo) AppError {
	return &appError{
		info:  info,
		cause: err,
	}
}
