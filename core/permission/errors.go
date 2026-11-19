package permission

import "errors"

var (
	ErrDenied = errors.New(
		"permission denied",
	)

	ErrApprovalRequired = errors.New(
		"approval required",
	)
)
