package util

type Result[T any] struct {
	Data []T `json:"data"`
}

func NewResult[T any](data []T) Result[T] {
	if data == nil {
		data = make([]T, 0)
	}
	return Result[T]{Data: data}
}
