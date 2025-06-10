package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pumpkinlog/backend/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	code := cmd.Execute(ctx)

	os.Exit(code)
}
