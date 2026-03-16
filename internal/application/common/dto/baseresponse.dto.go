package appdto

type BaseResponse[T any] struct {
	Data  *T
	Error *ErrorResponse
}

type ErrorResponse struct {
	Code    string
	Message string
	Detail  *string
}
