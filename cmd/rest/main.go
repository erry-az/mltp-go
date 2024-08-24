package main

import (
	"log/slog"
	"os"

	"github.com/erry-az/mltp-go/internal/app"
	"github.com/erry-az/mltp-go/internal/config"
)

func main() {
	cfg, err := config.Load("local")
	if err != nil {
		slog.Error("failed to load config", err)
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	err = app.NewRest(cfg)
	if err != nil {
		slog.Error("failed to start rest server", err)
	}
}
