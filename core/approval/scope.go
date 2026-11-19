package approval

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type scopePayload struct {
	TaskID string `json:"task_id"`

	SessionID string `json:"session_id"`

	ToolName string `json:"tool_name"`

	Arguments map[string]any `json:"arguments"`
}

func ScopeHash(
	taskID string,
	sessionID string,
	toolName string,
	arguments map[string]any,
) (string, error) {
	payload := scopePayload{
		TaskID: taskID,

		SessionID: sessionID,

		ToolName: toolName,

		Arguments: arguments,
	}

	data, err := json.Marshal(
		payload,
	)

	if err != nil {
		return "", fmt.Errorf(
			"marshal approval scope: %w",
			err,
		)
	}

	sum := sha256.Sum256(
		data,
	)

	return hex.EncodeToString(
		sum[:],
	), nil
}
