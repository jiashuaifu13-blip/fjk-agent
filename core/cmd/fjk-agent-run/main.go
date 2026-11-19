package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/fjk-dev/fjk-agent/core/agent"
	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/planner"
	fjkruntime "github.com/fjk-dev/fjk-agent/core/runtime"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	if cfg.DeepSeekAPIKey == "" {
		log.Fatal(
			"DEEPSEEK_API_KEY is not set",
		)
	}

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			nil,
		),
	)

	core, err :=
		fjkruntime.New(
			cfg,
			logger,
		)

	if err != nil {
		log.Fatal(err)
	}

	provider, err :=
		core.Models().Get(
			cfg.AIProvider,
		)

	if err != nil {
		log.Fatal(
			"get model provider: ",
			err,
		)
	}

	modelPlanner :=
		planner.NewModelPlanner(
			provider,
			core.Tools(),
		)

	runner :=
		agent.NewRunner(
			core.Tasks(),
			modelPlanner,
			core.Executor(),
		)

	run, err :=
		runner.Run(
			ctx,
			"session_agent_check",
			"请调用 system.echo 工具，把文本「FJK Agent end-to-end successful」原样返回。",
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"TASK ID:",
		run.TaskID,
	)

	fmt.Println(
		"RUN STATUS:",
		run.Status,
	)

	fmt.Println(
		"OUTPUT:",
		run.Output,
	)
}
