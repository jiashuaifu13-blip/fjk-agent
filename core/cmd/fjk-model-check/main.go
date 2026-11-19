package main

import (
	"context"
	"fmt"
	"log"

	deepseekadapter "github.com/fjk-dev/fjk-agent/core/adapters/ai/deepseek"
	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/model"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	if cfg.DeepSeekAPIKey == "" {
		log.Fatal(
			"DEEPSEEK_API_KEY is not set",
		)
	}

	provider, err := deepseekadapter.New(
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

	response, err := provider.Generate(
		ctx,
		model.Request{
			Messages: []model.Message{
				{
					Role: model.RoleSystem,
					Content: "You are the AI model used by FJK Agent Core. " +
						"Reply concisely.",
				},
				{
					Role: model.RoleUser,
					Content: "Reply exactly with: " +
						"FJK Agent model connection successful",
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"PROVIDER:",
		provider.Name(),
	)

	fmt.Println(
		"RESPONSE:",
		response.Content,
	)
}
