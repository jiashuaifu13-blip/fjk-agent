package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/fjk-dev/fjk-agent/core/executor"
	"github.com/fjk-dev/fjk-agent/core/planner"
	"github.com/fjk-dev/fjk-agent/core/task"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Runner struct {
	tasks *task.Manager

	planner planner.Planner

	executor *executor.Executor
}

func NewRunner(
	tasks *task.Manager,
	planner planner.Planner,
	executor *executor.Executor,
) *Runner {
	return &Runner{
		tasks: tasks,

		planner: planner,

		executor: executor,
	}
}

func (r *Runner) Run(
	ctx context.Context,
	sessionID string,
	input string,
) (*Run, error) {
	taskItem, err :=
		r.tasks.Create(
			ctx,
			sessionID,
			input,
			input,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"create task: %w",
			err,
		)
	}

	if _, err :=
		r.tasks.Transition(
			ctx,
			taskItem.ID,
			task.StatusQueued,
		); err != nil {
		return nil, err
	}

	if _, err :=
		r.tasks.Transition(
			ctx,
			taskItem.ID,
			task.StatusRunning,
		); err != nil {
		return nil, err
	}

	run := &Run{
		TaskID: taskItem.ID,

		SessionID: sessionID,

		Status: RunRunning,

		Input: input,
	}

	plan, err :=
		r.planner.Plan(
			ctx,
			planner.Request{
				TaskID: taskItem.ID,

				SessionID: sessionID,

				Input: input,
			},
		)

	if err != nil {
		_, _ =
			r.tasks.Transition(
				ctx,
				taskItem.ID,
				task.StatusFailed,
			)

		run.Status = RunFailed

		return run, fmt.Errorf(
			"plan task: %w",
			err,
		)
	}

	var output strings.Builder

	for _, step := range plan.Steps {

		switch step.Type {
		case planner.StepTypeReply:
			output.WriteString(
				step.Content,
			)

		case planner.StepTypeTool:
			result, err :=
				r.executor.ExecuteTool(
					ctx,
					tool.Request{
						TaskID: taskItem.ID,

						SessionID: sessionID,

						ToolName: step.ToolName,

						Arguments: step.Arguments,
					},
				)

			if err != nil {
				return run, err
			}

			fmt.Fprint(
				&output,
				result.Output,
			)

		default:
			return run, fmt.Errorf(
				"unsupported plan step type: %s",
				step.Type,
			)
		}
	}

	run.Output =
		output.String()

	run.Status =
		RunCompleted

	if _, err :=
		r.tasks.Transition(
			ctx,
			taskItem.ID,
			task.StatusCompleted,
		); err != nil {
		return run, err
	}

	return run, nil
}
