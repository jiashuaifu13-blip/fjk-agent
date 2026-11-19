package task

import (
	"context"
	"time"
)

type IDAllocator interface {
	NextTaskID(
		context.Context,
		time.Time,
	) (string, error)
}
