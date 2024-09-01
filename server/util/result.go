package util

type Result[T any] struct {
	Data []T `json:"data"`
}

func NewResult[T any](data []T) Result[T] {
	return Result[T]{Data: data}
}
