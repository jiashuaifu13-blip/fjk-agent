package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/executor"
	"github.com/fjk-dev/fjk-agent/core/permission"
	"github.com/fjk-dev/fjk-agent/core/planner"
	"github.com/fjk-dev/fjk-agent/core/task"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Runner struct {
	tasks *task.Manager

	planner planner.Planner

	executor *executor.Executor

	events *event.Bus
}

func NewRunner(
	tasks *task.Manager,
	planner planner.Planner,
	executor *executor.Executor,
) *Runner {
	return &Runner{
		tasks:    tasks,
		planner:  planner,
		executor: executor,
	}
}

func NewRunnerWithEvents(
	tasks *task.Manager,
	planner planner.Planner,
	executor *executor.Executor,
	events *event.Bus,
) *Runner {
	return &Runner{
		tasks:    tasks,
		planner:  planner,
		executor: executor,
		events:   events,
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
		TaskID:    taskItem.ID,
		SessionID: sessionID,
		Status:    RunRunning,
		Input:     input,
	}

	r.publish(
		ctx,
		event.TypeAgentRunStarted,
		run,
		map[string]any{
			"input": input,
		},
	)

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
		return r.fail(
			ctx,
			run,
			fmt.Errorf(
				"plan task: %w",
				err,
			),
		)
	}

	r.publish(
		ctx,
		event.TypeAgentPlanCreated,
		run,
		map[string]any{
			"step_count": len(plan.Steps),
		},
	)

	var output strings.Builder

	for index, step := range plan.Steps {

		r.publish(
			ctx,
			event.TypeAgentStepStarted,
			run,
			map[string]any{
				"step_id": step.ID,

				"step_index": index,

				"step_type": step.Type,

				"title": step.Title,

				"tool_name": step.ToolName,
			},
		)

		switch step.Type {
		case planner.StepTypeReply:
			output.WriteString(
				step.Content,
			)

		case planner.StepTypeTool:
			if r.executor == nil {
				return r.fail(
					ctx,
					run,
					fmt.Errorf(
						"tool executor is not available",
					),
				)
			}

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
				if errors.Is(
					err,
					permission.ErrApprovalRequired,
				) {
					run.Status =
						RunWaitingApproval

					r.publish(
						ctx,
						event.TypeAgentRunWaitingApproval,
						run,
						map[string]any{
							"step_id": step.ID,

							"tool_name": step.ToolName,

							"error": err.Error(),
						},
					)

					return run, err
				}

				return r.fail(
					ctx,
					run,
					err,
				)
			}

			fmt.Fprint(
				&output,
				result.Output,
			)

		default:
			return r.fail(
				ctx,
				run,
				fmt.Errorf(
					"unsupported plan step type: %s",
					step.Type,
				),
			)
		}

		r.publish(
			ctx,
			event.TypeAgentStepCompleted,
			run,
			map[string]any{
				"step_id": step.ID,

				"step_index": index,

				"step_type": step.Type,

				"tool_name": step.ToolName,
			},
		)
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

		return r.fail(
			ctx,
			run,
			err,
		)
	}

	r.publish(
		ctx,
		event.TypeAgentRunCompleted,
		run,
		map[string]any{
			"output": run.Output,
		},
	)

	return run, nil
}

func (r *Runner) fail(
	ctx context.Context,
	run *Run,
	runErr error,
) (*Run, error) {
	run.Status =
		RunFailed

	current, err :=
		r.tasks.Get(
			ctx,
			run.TaskID,
		)

	if err == nil &&
		current.Status == task.StatusRunning {
		_, _ =
			r.tasks.Transition(
				ctx,
				run.TaskID,
				task.StatusFailed,
			)
	}

	r.publish(
		ctx,
		event.TypeAgentRunFailed,
		run,
		map[string]any{
			"error": runErr.Error(),
		},
	)

	return run, runErr
}

func (r *Runner) publish(
	ctx context.Context,
	eventType event.Type,
	run *Run,
	payload map[string]any,
) {
	if r.events == nil {
		return
	}

	r.events.Publish(
		ctx,
		event.New(
			eventType,
			event.Metadata{
				TaskID: run.TaskID,

				SessionID: run.SessionID,
			},
			payload,
		),
	)
}
