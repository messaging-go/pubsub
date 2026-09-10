package integration_template

import (
	"context"
)

type Handler[T any] struct {
}

func (h Handler[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	panic("implement me")
}

func New[T any]() *Handler[T] {
	return &Handler[T]{}
}
