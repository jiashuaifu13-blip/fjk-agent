package executor

import (
	"context"
	"fmt"

	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type Executor struct {
	tools  *tool.Registry
	events *event.Bus
}

func New(
	tools *tool.Registry,
	events *event.Bus,
) *Executor {
	return &Executor{
		tools:  tools,
		events: events,
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
