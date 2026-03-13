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

type Option func(*appError)

func WithDetail(detail string) Option {
	return func(e *appError) {
		e.info.Detail = &detail
	}
}

func WithCause(err error) Option {
	return func(e *appError) {
		e.cause = err
	}
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

func New(info ErrorInfo, opts ...Option) AppError {
	e := &appError{
		info: info,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func Wrap(err error, info ErrorInfo, opts ...Option) AppError {
	e := &appError{
		info: info,
	}

	for _, opt := range opts {
		opt(e)
	}

	return &appError{
		info:  info,
		cause: err,
	}
}
