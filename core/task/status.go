package task

type Status string

const (
	StatusCreated         Status = "CREATED"
	StatusQueued          Status = "QUEUED"
	StatusRunning         Status = "RUNNING"
	StatusWaitingApproval Status = "WAITING_APPROVAL"
	StatusPaused          Status = "PAUSED"
	StatusCompleted       Status = "COMPLETED"
	StatusFailed          Status = "FAILED"
	StatusCancelled       Status = "CANCELLED"
)

func (s Status) Valid() bool {
	switch s {
	case
		StatusCreated,
		StatusQueued,
		StatusRunning,
		StatusWaitingApproval,
		StatusPaused,
		StatusCompleted,
		StatusFailed,
		StatusCancelled:
		return true
	default:
		return false
	}
}
