package task

import (
	"context"
	"fmt"
	"time"
)

func allocateID(
	ctx context.Context,
	repository Repository,
	now time.Time,
) (string, error) {
	if allocator, ok :=
		repository.(IDAllocator); ok {

		id, err :=
			allocator.NextTaskID(
				ctx,
				now,
			)

		if err != nil {
			return "", fmt.Errorf(
				"allocate persistent task id: %w",
				err,
			)
		}

		return id, nil
	}

	return NewID(), nil
}
