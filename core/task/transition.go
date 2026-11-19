package task

import (
	"fmt"
	"time"
)

func (t *Task) Transition(next Status) error {
	if !next.Valid() {
		return fmt.Errorf("invalid task status: %s", next)
	}

	if !canTransition(t.Status, next) {
		return fmt.Errorf(
			"invalid task transition: %s -> %s",
			t.Status,
			next,
		)
	}

	now := time.Now().UTC()

	if next == StatusRunning && t.StartedAt == nil {
		t.StartedAt = &now
	}

	if next == StatusCompleted ||
		next == StatusFailed ||
		next == StatusCancelled {
		t.CompletedAt = &now
	}

	t.Status = next
	t.UpdatedAt = now

	return nil
}

func canTransition(current, next Status) bool {
	switch current {
	case StatusCreated:
		return next == StatusQueued ||
			next == StatusCancelled

	case StatusQueued:
		return next == StatusRunning ||
			next == StatusCancelled

	case StatusRunning:
		return next == StatusWaitingApproval ||
			next == StatusPaused ||
			next == StatusCompleted ||
			next == StatusFailed ||
			next == StatusCancelled

	case StatusWaitingApproval:
		return next == StatusRunning ||
			next == StatusCancelled ||
			next == StatusFailed

	case StatusPaused:
		return next == StatusRunning ||
			next == StatusCancelled

	case StatusCompleted,
		StatusFailed,
		StatusCancelled:
		return false

	default:
		return false
	}
}
