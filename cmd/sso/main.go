package main

import (
	"log/slog"
	"os"

	"github.com/m1al04949/sso-gRPC/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Config initialize
	cfg := config.MustLoad()

	// Logger initialize
	log := setupLogger(cfg.Env)

	log.Info("starting application",
		slog.Any("cfg", cfg))

	// TODO: Инициализация приложения

	// TODO: Запуск gRPC-сервис
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
