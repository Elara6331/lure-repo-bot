package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.arsenm.dev/lure-repo-bot/internal/queue"
	"go.arsenm.dev/lure-repo-bot/internal/spdx"
	"go.arsenm.dev/lure-repo-bot/internal/types"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	spdx.StartUpdater(ctx)

	jobQueue := queue.New[*types.PullRequestPayload]()

	startWebhookWorkers(ctx, jobQueue)

	addr := ":8080"
	if os.Getenv("LURE_BOT_ADDR") != "" {
		addr = os.Getenv("LURE_BOT_ADDR")
	}

	serveWebhook(ctx, addr, jobQueue)
}
