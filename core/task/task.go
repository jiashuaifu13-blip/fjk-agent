package task

import "time"

type Task struct {
	ID             string     `json:"task_id"`
	SessionID      string     `json:"session_id,omitempty"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	Status         Status     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	FailureMessage string     `json:"failure_message,omitempty"`
}

func New(id, sessionID, title, description string) *Task {
	now := time.Now().UTC()

	return &Task{
		ID:          id,
		SessionID:   sessionID,
		Title:       title,
		Description: description,
		Status:      StatusCreated,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
