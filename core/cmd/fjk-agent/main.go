package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fjk-dev/fjk-agent/core/config"
	"github.com/fjk-dev/fjk-agent/core/internal/buildinfo"
	"github.com/fjk-dev/fjk-agent/core/runtime"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	cfg := config.Load()

	logger.Info(
		"starting application",
		"name",
		buildinfo.Name,
		"version",
		buildinfo.Version,
		"environment",
		cfg.Environment,
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()

	core, err := runtime.New(
		cfg,
		logger,
	)

	if err != nil {
		logger.Error(
			"failed to initialize FJK Agent Core",
			"error",
			err,
		)

		os.Exit(1)
	}

	if err := core.Start(ctx); err != nil {
		logger.Error(
			"FJK Agent Core exited with error",
			"error",
			err,
		)

		os.Exit(1)
	}
}
