package event

import (
	"time"
)

type Type string

const (
	TypeCoreStarted  Type = "core.started"
	TypeCoreStopping Type = "core.stopping"
	TypeCoreStopped  Type = "core.stopped"
)

type Event struct {
	Type      Type           `json:"event_type"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload,omitempty"`
}

func New(eventType Type, payload map[string]any) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}
