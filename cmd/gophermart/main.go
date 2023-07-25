package main

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/server"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
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

	conn, err := database.NewConnection("pgx", cfg.GetStorageDSN())
	if err != nil {
		logger.Error("failed to create database connection", sl.Err(err))
		os.Exit(1)
	}
	defer func() {
		dbInstance, _ := conn.DB()
		_ = dbInstance.Close()
	}()
	logger.Info("db connection created")

	err = database.RunMigrations(cfg.GetStorageDSN())
	if err != nil {
		logger.Error("failed to run migrations", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("migrations started")

	storage := database.NewStorage(conn)

	router := server.NewRouter(logger, storage, &cfg).Serve()
	http.ListenAndServe(cfg.GetHTTPServerAddress(), router)

	// TODO: swagger
}
