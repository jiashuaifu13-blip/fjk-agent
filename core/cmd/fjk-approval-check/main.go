package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fjk-dev/fjk-agent/core/approval"
	"github.com/fjk-dev/fjk-agent/core/config"
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
		approval.NewSQLiteRepository(
			store.DB(),
		)

	manager :=
		approval.NewManager(
			repository,
		)

	arguments := map[string]any{
		"path":    "example.txt",
		"content": "FJK Agent approval test",
	}

	item, err :=
		manager.Request(
			ctx,
			"task_approval_test",
			"session_approval_test",
			"filesystem.write",
			arguments,
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"APPROVAL:",
		item.ID,
		item.Status,
	)

	if err := manager.Approve(
		ctx,
		item.ID,
	); err != nil {
		log.Fatal(err)
	}

	approved, err :=
		manager.Get(
			ctx,
			item.ID,
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"STATUS AFTER APPROVE:",
		approved.Status,
	)

	if err :=
		manager.ValidateAndConsume(
			ctx,
			item.ID,
			"task_approval_test",
			"session_approval_test",
			"filesystem.write",
			arguments,
		); err != nil {
		log.Fatal(err)
	}

	consumed, err :=
		manager.Get(
			ctx,
			item.ID,
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"STATUS AFTER CONSUME:",
		consumed.Status,
	)

	fmt.Println(
		"APPROVAL PERSISTENCE CHECK PASSED",
	)
}
