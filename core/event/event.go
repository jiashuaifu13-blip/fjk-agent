package event

import "time"

const ProtocolVersion = "1"

type Type string

const (
	TypeCoreStarted  Type = "core.started"
	TypeCoreStopping Type = "core.stopping"
	TypeCoreStopped  Type = "core.stopped"

	TypeTaskCreated         Type = "task.created"
	TypeTaskQueued          Type = "task.queued"
	TypeTaskStarted         Type = "task.started"
	TypeTaskProgress        Type = "task.progress"
	TypeTaskWaitingApproval Type = "task.waiting_approval"
	TypeTaskPaused          Type = "task.paused"
	TypeTaskCompleted       Type = "task.completed"
	TypeTaskFailed          Type = "task.failed"
	TypeTaskCancelled       Type = "task.cancelled"

	TypeToolCalled Type = "task.tool.called"
	TypeToolResult Type = "task.tool.result"

	TypeConversationStarted          Type = "conversation.started"
	TypeConversationMessageCreated   Type = "conversation.message.created"
	TypeConversationMessageDelta     Type = "conversation.message.delta"
	TypeConversationMessageCompleted Type = "conversation.message.completed"
	TypeConversationFailed           Type = "conversation.failed"

	TypePetStateChanged   Type = "pet.state.changed"
	TypePetEmotionChanged Type = "pet.emotion.changed"
)

type Event struct {
	Version   string         `json:"version"`
	EventID   string         `json:"event_id"`
	EventType Type           `json:"event_type"`
	TaskID    string         `json:"task_id,omitempty"`
	SessionID string         `json:"session_id,omitempty"`
	DeviceID  string         `json:"device_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload,omitempty"`
}
