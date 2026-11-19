package main

import (
"context"
"fmt"
"log"

"github.com/fjk-dev/fjk-agent/core/config"
"github.com/fjk-dev/fjk-agent/core/storage/sqlite"
"github.com/fjk-dev/fjk-agent/core/task"
)

func main() {
ctx := context.Background()

cfg := config.Load()

store, err := sqlite.Open(cfg.DatabasePath)
if err != nil {
log.Fatal(err)
}
defer store.Close()

repository := task.NewSQLiteRepository(store.DB())

items, err := repository.List(ctx)
if err != nil {
log.Fatal(err)
}

fmt.Println("TASK COUNT:", len(items))

for _, item := range items {
fmt.Println(
item.ID,
item.Status,
item.Title,
)
}
}
