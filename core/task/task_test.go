package task

import "testing"

func TestTaskLifecycle(t *testing.T) {
	item := New(
		NewID(),
		"session_test",
		"Test Task",
		"Task lifecycle test",
	)

	if item.Status != StatusCreated {
		t.Fatalf(
			"expected CREATED, got %s",
			item.Status,
		)
	}

	transitions := []Status{
		StatusQueued,
		StatusRunning,
		StatusWaitingApproval,
		StatusRunning,
		StatusCompleted,
	}

	for _, next := range transitions {
		if err := item.Transition(next); err != nil {
			t.Fatalf(
				"transition to %s failed: %v",
				next,
				err,
			)
		}
	}

	if item.Status != StatusCompleted {
		t.Fatalf(
			"expected COMPLETED, got %s",
			item.Status,
		)
	}

	if item.StartedAt == nil {
		t.Fatal("StartedAt was not set")
	}

	if item.CompletedAt == nil {
		t.Fatal("CompletedAt was not set")
	}
}

func TestInvalidTaskTransition(t *testing.T) {
	item := New(
		NewID(),
		"",
		"Invalid Transition",
		"",
	)

	err := item.Transition(StatusCompleted)

	if err == nil {
		t.Fatal(
			"expected invalid transition error",
		)
	}
}
