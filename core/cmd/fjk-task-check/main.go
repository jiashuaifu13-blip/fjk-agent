package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/event"
	"github.com/fjk-dev/fjk-agent/core/storage/sqlite"
	"github.com/fjk-dev/fjk-agent/core/task"
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

	if err :=
		store.Migrate(); err != nil {
		log.Fatal(err)
	}

	manager :=
		task.NewManager(
			task.NewSQLiteRepository(
				store.DB(),
			),
			event.NewBus(),
		)

	item, err :=
		manager.Create(
			ctx,
			"session_id_check",
			"Persistent Task ID Check",
			"Verify task IDs survive process restart",
		)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"TASK ID:",
		item.ID,
	)
}
