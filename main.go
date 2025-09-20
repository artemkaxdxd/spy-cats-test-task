package main

import (
	"backend/config"
	"backend/internal/app"
	"log/slog"
	"os"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("config fill", "err", err)
		os.Exit(1)
	}

	app.Run(cfg)
}
