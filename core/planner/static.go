package planner

import (
	"context"
	"fmt"
)

type StaticPlanner struct{}

func NewStaticPlanner() *StaticPlanner {
	return &StaticPlanner{}
}

func (p *StaticPlanner) Plan(
	ctx context.Context,
	request Request,
) (*Plan, error) {
	if request.TaskID == "" {
		return nil, fmt.Errorf(
			"task id is required",
		)
	}

	return &Plan{
		TaskID: request.TaskID,

		Steps: []Step{
			{
				ID: "step_001",

				Type: StepTypeReply,

				Title: "Respond to user",

				Content: request.Input,
			},
		},
	}, nil
}
