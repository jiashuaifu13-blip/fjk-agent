package executor

import (
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/permission"
)

type ApprovalRequiredError struct {
	ApprovalID string

	Reason string
}

func (e *ApprovalRequiredError) Error() string {
	if e.ApprovalID == "" {
		return fmt.Sprintf(
			"approval required: %s",
			e.Reason,
		)
	}

	return fmt.Sprintf(
		"approval required: %s (approval_id=%s)",
		e.Reason,
		e.ApprovalID,
	)
}

func (e *ApprovalRequiredError) Unwrap() error {
	return permission.ErrApprovalRequired
}
