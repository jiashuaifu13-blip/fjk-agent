package approval

import "errors"

var (
	ErrNotFound = errors.New(
		"approval not found",
	)

	ErrInvalidStatus = errors.New(
		"invalid approval status",
	)

	ErrNotPending = errors.New(
		"approval is not pending",
	)

	ErrNotApproved = errors.New(
		"approval is not approved",
	)

	ErrScopeMismatch = errors.New(
		"approval scope mismatch",
	)

	ErrAlreadyConsumed = errors.New(
		"approval already consumed",
	)
)
