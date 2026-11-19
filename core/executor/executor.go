package executor

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/permission"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Executor struct {
	tools       *tool.Registry
	events      *event.Bus
	permissions *permission.Engine
}

func New(
	tools *tool.Registry,
	events *event.Bus,
	permissions *permission.Engine,
) *Executor {
	return &Executor{
		tools:       tools,
		events:      events,
		permissions: permissions,
	}
}

func (e *Executor) ExecuteTool(
	ctx context.Context,
	request tool.Request,
) (*tool.Result, error) {
	item, err :=
		e.tools.Get(request.ToolName)

	if err != nil {
		return nil, err
	}

	definition := item.Definition()

	required :=
		permission.FromToolRisk(
			definition.Risk,
		)

	decision :=
		e.permissions.Check(
			permission.Request{
				TaskID:       request.TaskID,
				SessionID:    request.SessionID,
				ToolName:     request.ToolName,
				Required:     required,
				UserApproved: request.UserApproved,
			},
		)

	if !decision.Allowed {
		if decision.RequiresApproval {
			return nil, fmt.Errorf(
				"%w: %s",
				permission.ErrApprovalRequired,
				decision.Reason,
			)
		}

		return nil, fmt.Errorf(
			"%w: %s",
			permission.ErrDenied,
			decision.Reason,
		)
	}

	e.events.Publish(
		ctx,
		event.New(
			event.TypeToolCalled,
			event.Metadata{
				TaskID:    request.TaskID,
				SessionID: request.SessionID,
			},
			map[string]any{
				"tool_name": definition.Name,
				"risk":      definition.Risk,
				"arguments": request.Arguments,
			},
		),
	)

	result, err := item.Execute(
		ctx,
		request,
	)

	if err != nil {
		e.events.Publish(
			ctx,
			event.New(
				event.TypeToolResult,
				event.Metadata{
					TaskID:    request.TaskID,
					SessionID: request.SessionID,
				},
				map[string]any{
					"tool_name": definition.Name,
					"success":   false,
					"error":     err.Error(),
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
				TaskID:    request.TaskID,
				SessionID: request.SessionID,
			},
			map[string]any{
				"tool_name": definition.Name,
				"success":   true,
				"output":    result.Output,
			},
		),
	)

	return result, nil
}
