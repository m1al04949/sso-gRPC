package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1al04949/sso-gRPC/internal/app"
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

	// Initialize App
	appl := app.New(log, cfg.GRPC.Port, cfg.DB, cfg.JWT.TokenTTL)

	// gRPC Server Run
	go appl.GRPCSrv.MustRun()

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signalIn := <-stop
	log.Info("stopping application",
		slog.String("signal", signalIn.String()))

	// Storage stop
	appl.Storage.Close()

	// gRPC stop
	appl.GRPCSrv.Stop()

	log.Info("application stopped")
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
