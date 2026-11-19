package planner

type StepType string

const (
	StepTypeTool  StepType = "TOOL"
	StepTypeReply StepType = "REPLY"
)

type Step struct {
	ID string `json:"id"`

	Type StepType `json:"type"`

	Title string `json:"title,omitempty"`

	ToolName string `json:"tool_name,omitempty"`

	Arguments map[string]any `json:"arguments,omitempty"`

	Content string `json:"content,omitempty"`
}

type Plan struct {
	TaskID string `json:"task_id"`

	Steps []Step `json:"steps"`
}
