package appdto

type Pagination struct {
	Page  int
	Limit int
}

type PaginatedResponse[T any] struct {
	Items []T
	Total int
	Page  int
	Limit int
	Size  int
}
