package tool

import "errors"

var (
	ErrToolNotFound = errors.New(
		"tool not found",
	)

	ErrToolAlreadyRegistered = errors.New(
		"tool already registered",
	)

	ErrInvalidToolDefinition = errors.New(
		"invalid tool definition",
	)
)
