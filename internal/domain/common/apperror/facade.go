package errorbase

type AppError interface {
	error
	ErrorInfo() ErrorInfo
}
