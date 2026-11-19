package main

import (
	"context"
	"fmt"
	"log"

	deepseekadapter "github.com/fjk-dev/fjk-agent/core/adapters/ai/deepseek"
	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/conversation"
	"github.com/fjk-dev/fjk-agent/core/event"
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

	bus := event.NewBus()

	manager :=
		conversation.NewManager()

	service :=
		conversation.NewService(
			manager,
			provider,
			bus,
		)

	bus.Subscribe(
		event.TypeConversationMessageDelta,
		func(
			ctx context.Context,
			evt event.Event,
		) {
			if delta, ok :=
				evt.Payload["delta"].(string); ok {
				fmt.Print(delta)
			}
		},
	)

	bus.Subscribe(
		event.TypeConversationMessageCompleted,
		func(
			ctx context.Context,
			evt event.Event,
		) {
			fmt.Println()
			fmt.Println(
				"[conversation completed]",
			)
		},
	)

	conv :=
		manager.Create(
			ctx,
			"Streaming Test",
		)

	fmt.Println(
		"Conversation:",
		conv.ID,
	)

	fmt.Println("Assistant:")

	_, err = service.Chat(
		ctx,
		conv.ID,
		"请用一句中文说明你已经通过流式输出连接到了 FJK Agent。",
	)

	if err != nil {
		log.Fatal(err)
	}
}
