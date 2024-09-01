package util

type Result[T any] struct {
	data []T `json:"data"`
}

func NewResult[T any](data []T) Result[T] {
	return Result[T]{data: data}
}
