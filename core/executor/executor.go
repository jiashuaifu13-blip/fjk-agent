package executor

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/approval"
	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/permission"
	"github.com/fjk-dev/fjk-agent/core/task"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Executor struct {
	tools *tool.Registry

	events *event.Bus

	permissions *permission.Engine

	approvals *approval.Manager

	tasks *task.Manager
}

func New(
	tools *tool.Registry,
	events *event.Bus,
	permissions *permission.Engine,
	approvals *approval.Manager,
	tasks *task.Manager,
) *Executor {
	return &Executor{
		tools: tools,

		events: events,

		permissions: permissions,

		approvals: approvals,

		tasks: tasks,
	}
}

func (e *Executor) ExecuteTool(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	item, err :=
		e.tools.Get(
			request.ToolName,
		)

	if err != nil {
		return nil, err
	}

	definition :=
		item.Definition()

	required :=
		permission.FromToolRisk(
			definition.Risk,
		)

	decision :=
		e.permissions.Check(
			permission.Request{
				TaskID: request.TaskID,

				SessionID: request.SessionID,

				ToolName: request.ToolName,

				Required: required,
			},
		)

	if !decision.Allowed {
		if !decision.RequiresApproval {
			return nil, fmt.Errorf(
				"%w: %s",
				permission.ErrDenied,
				decision.Reason,
			)
		}

		if request.ApprovalID == "" {
			approvalItem, err :=
				e.approvals.Request(
					ctx,
					request.TaskID,
					request.SessionID,
					request.ToolName,
					request.Arguments,
				)

			if err != nil {
				return nil, fmt.Errorf(
					"create approval: %w",
					err,
				)
			}

			if err :=
				e.waitForApproval(
					ctx,
					request.TaskID,
				); err != nil {
				return nil, err
			}

			return nil,
				&ApprovalRequiredError{
					ApprovalID: approvalItem.ID,
					Reason:     decision.Reason,
				}
		}

		if e.approvals == nil {
			return nil, fmt.Errorf(
				"%w: approval manager unavailable",
				permission.ErrDenied,
			)
		}

		if err :=
			e.approvals.ValidateAndConsume(
				ctx,
				request.ApprovalID,
				request.TaskID,
				request.SessionID,
				request.ToolName,
				request.Arguments,
			); err != nil {

			return nil, fmt.Errorf(
				"validate approval: %w",
				err,
			)
		}

		if err :=
			e.resumeAfterApproval(
				ctx,
				request.TaskID,
			); err != nil {
			return nil, err
		}
	}

	e.events.Publish(
		ctx,
		event.New(
			event.TypeToolCalled,
			event.Metadata{
				TaskID: request.TaskID,

				SessionID: request.SessionID,
			},
			map[string]any{
				"tool_name": definition.Name,

				"risk": definition.Risk,

				"arguments": request.Arguments,

				"approval_id": request.ApprovalID,
			},
		),
	)

	result, err :=
		item.Execute(
			ctx,
			request,
		)

	if err != nil {
		e.events.Publish(
			ctx,
			event.New(
				event.TypeToolResult,
				event.Metadata{
					TaskID: request.TaskID,

					SessionID: request.SessionID,
				},
				map[string]any{
					"tool_name": definition.Name,

					"success": false,

					"error": err.Error(),
				},
			),
		)

		return nil, fmt.Errorf(
			"execute tool %s: %w",
			definition.Name,
			err,
		)
	}

	e.events.Publish(
		ctx,
		event.New(
			event.TypeToolResult,
			event.Metadata{
				TaskID: request.TaskID,

				SessionID: request.SessionID,
			},
			map[string]any{
				"tool_name": definition.Name,

				"success": true,

				"output": result.Output,
			},
		),
	)

	return result, nil
}

func (e *Executor) waitForApproval(
	ctx context.Context,
	taskID string,
) error {
	if taskID == "" ||
		e.tasks == nil {
		return nil
	}

	current, err :=
		e.tasks.Get(
			ctx,
			taskID,
		)

	if err != nil {
		return fmt.Errorf(
			"get task before approval: %w",
			err,
		)
	}

	if current.Status != task.StatusRunning {
		return nil
	}

	if _, err :=
		e.tasks.Transition(
			ctx,
			taskID,
			task.StatusWaitingApproval,
		); err != nil {

		return fmt.Errorf(
			"transition task to waiting approval: %w",
			err,
		)
	}

	return nil
}

func (e *Executor) resumeAfterApproval(
	ctx context.Context,
	taskID string,
) error {
	if taskID == "" ||
		e.tasks == nil {
		return nil
	}

	current, err :=
		e.tasks.Get(
			ctx,
			taskID,
		)

	if err != nil {
		return fmt.Errorf(
			"get task after approval: %w",
			err,
		)
	}

	if current.Status !=
		task.StatusWaitingApproval {
		return nil
	}

	if _, err :=
		e.tasks.Transition(
			ctx,
			taskID,
			task.StatusRunning,
		); err != nil {

		return fmt.Errorf(
			"resume task after approval: %w",
			err,
		)
	}

	return nil
}
