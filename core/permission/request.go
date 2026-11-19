package permission

type Request struct {
	TaskID       string
	SessionID    string
	ToolName     string
	Required     Level
	UserApproved bool
}
