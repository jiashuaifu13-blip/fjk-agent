package approval

import "time"

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
	StatusConsumed Status = "CONSUMED"
)

type Approval struct {
	ID string

	TaskID string

	SessionID string

	ToolName string

	ScopeHash string

	Status Status

	CreatedAt time.Time

	DecidedAt *time.Time

	ConsumedAt *time.Time
}
