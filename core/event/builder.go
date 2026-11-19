package event

import "time"

type Metadata struct {
	TaskID    string
	SessionID string
	DeviceID  string
}

func New(
	eventType Type,
	meta Metadata,
	payload map[string]any,
) Event {
	return Event{
		Version:   ProtocolVersion,
		EventID:   NewID(),
		EventType: eventType,
		TaskID:    meta.TaskID,
		SessionID: meta.SessionID,
		DeviceID:  meta.DeviceID,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}
