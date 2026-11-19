package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/conversation"
	"github.com/fjk-dev/fjk-agent/core/storage/sqlite"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	store, err :=
		sqlite.Open(
			cfg.DatabasePath,
		)

	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	if err := store.Migrate(); err != nil {
		log.Fatal(err)
	}

	repository :=
		conversation.NewSQLiteRepository(
			store.DB(),
		)

	manager :=
		conversation.NewManagerWithRepository(
			repository,
		)

	conv :=
		manager.Create(
			ctx,
			"SQLite Conversation Test",
		)

	if conv == nil {
		log.Fatal(
			"failed to create conversation",
		)
	}

	manager.AddMessage(
		ctx,
		conv.ID,
		conversation.RoleUser,
		"Hello FJK Agent",
	)

	manager.AddMessage(
		ctx,
		conv.ID,
		conversation.RoleAssistant,
		"Hello from SQLite",
	)

	stored, err :=
		manager.Get(
			ctx,
			conv.ID,
		)

	if err != nil {
		log.Fatal(err)
	}

	messages :=
		manager.Messages(
			ctx,
			conv.ID,
		)

	fmt.Println(
		"CONVERSATION:",
		stored.ID,
		stored.Title,
	)

	fmt.Println(
		"MESSAGE COUNT:",
		len(messages),
	)

	for _, message := range messages {
		fmt.Println(
			message.Role,
			":",
			message.Content,
		)
	}

	if len(messages) != 2 {
		log.Fatal(
			"conversation persistence check failed",
		)
	}

	fmt.Println(
		"CONVERSATION PERSISTENCE CHECK PASSED",
	)
}
