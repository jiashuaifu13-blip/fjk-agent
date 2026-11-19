package event

import "testing"

func TestNewEvent(t *testing.T) {
	evt := New(
		TypeTaskCreated,
		Metadata{
			TaskID:    "task_test_001",
			SessionID: "session_test",
			DeviceID:  "device_test",
		},
		map[string]any{
			"title": "test",
		},
	)

	if evt.Version != ProtocolVersion {
		t.Fatalf(
			"unexpected protocol version: %s",
			evt.Version,
		)
	}

	if evt.EventID == "" {
		t.Fatal("event id is empty")
	}

	if evt.EventType != TypeTaskCreated {
		t.Fatalf(
			"unexpected event type: %s",
			evt.EventType,
		)
	}

	if evt.TaskID != "task_test_001" {
		t.Fatalf(
			"unexpected task id: %s",
			evt.TaskID,
		)
	}
}
