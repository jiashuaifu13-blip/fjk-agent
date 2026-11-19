package task

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(context.Context, *Task) error
	Update(context.Context, *Task) error
	Get(context.Context, string) (*Task, error)
	List(context.Context) ([]*Task, error)
}
