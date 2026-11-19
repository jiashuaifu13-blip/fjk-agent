package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	deepseekadapter "github.com/fjk-dev/fjk-agent/core/adapters/ai/deepseek"
	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/planner"
	"github.com/fjk-dev/fjk-agent/core/tool"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	if cfg.DeepSeekAPIKey == "" {
		log.Fatal(
			"DEEPSEEK_API_KEY is not set",
		)
	}

	provider, err :=
		deepseekadapter.New(
			ctx,
			deepseekadapter.Config{
				APIKey:  cfg.DeepSeekAPIKey,
				BaseURL: cfg.DeepSeekBaseURL,
				Model:   cfg.AIModel,
			},
		)

	if err != nil {
		log.Fatal(err)
	}

	modelPlanner :=
		planner.NewModelPlanner(
			provider,
			tool.NewRegistry(),
		)

	plan, err :=
		modelPlanner.Plan(
			ctx,
			planner.Request{
				TaskID:    "task_planner_check",
				SessionID: "session_planner_check",
				Input:     "请用一句中文回复：FJK Agent Planner 已连接成功。",
			},
		)

	if err != nil {
		log.Fatal(err)
	}

	data, err :=
		json.MarshalIndent(
			plan,
			"",
			"  ",
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		string(data),
	)
}
