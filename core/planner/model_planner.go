package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fjk-dev/fjk-agent/core/model"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

type ModelPlanner struct {
	provider model.Provider

	tools *tool.Registry
}

func NewModelPlanner(
	provider model.Provider,
	tools *tool.Registry,
) *ModelPlanner {
	return &ModelPlanner{
		provider: provider,
		tools:    tools,
	}
}

func (p *ModelPlanner) Plan(
	ctx context.Context,
	request Request,
) (*Plan, error) {
	if request.TaskID == "" {
		return nil, fmt.Errorf(
			"task id is required",
		)
	}

	if p.provider == nil {
		return nil, fmt.Errorf(
			"model provider is required",
		)
	}

	toolDefinitions := []tool.Definition{}

	if p.tools != nil {
		toolDefinitions =
			p.tools.Definitions()
	}

	toolJSON, err :=
		json.Marshal(
			toolDefinitions,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"marshal tool definitions: %w",
			err,
		)
	}

	systemPrompt := fmt.Sprintf(
		`You are the planner for FJK Agent.

Your job is ONLY to create a plan.
You must never claim that a tool has already been executed.

Available tools:
%s

Return ONLY valid JSON using this schema:

{
  "task_id": "string",
  "steps": [
    {
      "id": "step_001",
      "type": "TOOL or REPLY",
      "title": "short description",
      "tool_name": "required for TOOL",
      "arguments": {},
      "content": "required for REPLY"
    }
  ]
}

Rules:
- Use only tools from the available tools list.
- If no tool is required, use one REPLY step.
- TOOL arguments must be JSON objects.
- Do not include markdown fences.
- Do not include explanations outside JSON.
- task_id must exactly equal %q.`,
		string(toolJSON),
		request.TaskID,
	)

	response, err :=
		p.provider.Generate(
			ctx,
			model.Request{
				Messages: []model.Message{
					{
						Role:    model.RoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    model.RoleUser,
						Content: request.Input,
					},
				},
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"generate plan: %w",
			err,
		)
	}

	raw :=
		extractJSONObject(
			response.Content,
		)

	var plan Plan

	if err := json.Unmarshal(
		[]byte(raw),
		&plan,
	); err != nil {
		return nil, fmt.Errorf(
			"decode model plan: %w",
			err,
		)
	}

	if err :=
		validatePlan(
			request,
			&plan,
			p.tools,
		); err != nil {
		return nil, err
	}

	return &plan, nil
}

func extractJSONObject(
	content string,
) string {
	content = strings.TrimSpace(
		content,
	)

	start := strings.Index(
		content,
		"{",
	)

	end := strings.LastIndex(
		content,
		"}",
	)

	if start >= 0 &&
		end >= start {
		return content[start : end+1]
	}

	return content
}

func validatePlan(
	request Request,
	plan *Plan,
	tools *tool.Registry,
) error {
	if plan.TaskID != request.TaskID {
		return fmt.Errorf(
			"plan task id mismatch: expected %s, got %s",
			request.TaskID,
			plan.TaskID,
		)
	}

	if len(plan.Steps) == 0 {
		return fmt.Errorf(
			"plan contains no steps",
		)
	}

	for index, step := range plan.Steps {

		if step.ID == "" {
			return fmt.Errorf(
				"plan step %d has no id",
				index,
			)
		}

		switch step.Type {
		case StepTypeReply:
			if step.Content == "" {
				return fmt.Errorf(
					"reply step %s has no content",
					step.ID,
				)
			}

		case StepTypeTool:
			if step.ToolName == "" {
				return fmt.Errorf(
					"tool step %s has no tool name",
					step.ID,
				)
			}

			if tools == nil {
				return fmt.Errorf(
					"tool step %s cannot be validated without registry",
					step.ID,
				)
			}

			if _, err :=
				tools.Get(
					step.ToolName,
				); err != nil {

				return fmt.Errorf(
					"plan step %s references unavailable tool %s: %w",
					step.ID,
					step.ToolName,
					err,
				)
			}

		default:
			return fmt.Errorf(
				"unsupported plan step type %q",
				step.Type,
			)
		}
	}

	return nil
}
