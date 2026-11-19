package agent

type RunStatus string

const (
	RunCreated         RunStatus = "CREATED"
	RunRunning         RunStatus = "RUNNING"
	RunWaitingApproval RunStatus = "WAITING_APPROVAL"
	RunCompleted       RunStatus = "COMPLETED"
	RunFailed          RunStatus = "FAILED"
)

type Run struct {
	TaskID string

	SessionID string

	Status RunStatus

	Input string

	Output string
}
