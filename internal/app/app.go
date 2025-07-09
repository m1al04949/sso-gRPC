package app

import (
	"log/slog"
	"time"

	"github.com/m1al04949/sso-gRPC/internal/app/grpcapp"
	"github.com/m1al04949/sso-gRPC/internal/config"
	"github.com/m1al04949/sso-gRPC/internal/services/auth"
	"github.com/m1al04949/sso-gRPC/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *sqlite.Storage
}

func New(log *slog.Logger, grpcPort int, dbCfg config.DBConfig, tokenTTL time.Duration) *App {
	// Init storage
	storage, err := sqlite.New(log, dbCfg)
	if err != nil {
		panic(err)
	}

	// Init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	// Init app
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
		Storage: storage,
	}
}
