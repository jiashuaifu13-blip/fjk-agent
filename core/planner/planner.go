package planner

import "context"

type Request struct {
	TaskID string

	SessionID string

	Input string
}

type Planner interface {
	Plan(
		context.Context,
		Request,
	) (*Plan, error)
}
