package tool

import "context"

type RiskLevel string

const (
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

type Definition struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Risk        RiskLevel `json:"risk"`
}

type Request struct {
	TaskID    string         `json:"task_id,omitempty"`
	SessionID string         `json:"session_id,omitempty"`
	ToolName  string         `json:"tool_name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type Result struct {
	ToolName string `json:"tool_name"`
	Output   any    `json:"output,omitempty"`
}

type Tool interface {
	Definition() Definition

	Execute(
		context.Context,
		Request,
	) (*Result, error)
}
