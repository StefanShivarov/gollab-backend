package common

type PaginationQuery struct {
	Page int `validate:"gte=1"`
	Size int `validate:"gte=1,lte=100"`
}

type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
