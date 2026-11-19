package echo

import (
	"context"
	"testing"

	"github.com/fjk-dev/fjk-agent/core/tool"
)

func TestEcho(t *testing.T) {
	item := New()

	result, err :=
		item.Execute(
			context.Background(),
			tool.Request{
				ToolName: "system.echo",

				Arguments: map[string]any{
					"text": "hello",
				},
			},
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.Output != "hello" {
		t.Fatalf(
			"unexpected output: %v",
			result.Output,
		)
	}
}
