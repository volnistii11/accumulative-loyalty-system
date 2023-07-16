package main

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/storage"
	"github.com/volnistii11/accumulative-loyalty-system/internal/storage/postgres"
	"golang.org/x/exp/slog"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.New()
	err := cfg.Parse()
	if err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	logger.Info("init cfg and logger completed")

	db, err := storage.NewConnection("pgx", cfg.GetStorageDSN())
	if err != nil {
		logger.Error("failed to create database connection", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("db connection created")

	err = postgres.RunMigrations(cfg.GetStorageDSN())
	if err != nil {
		logger.Error("failed to run migrations", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("migrations started")

	// TODO: init router: chi, render

	// TODO: run server

	// TODO: swagger
}
