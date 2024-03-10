package useCases

type PagedResult[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}
