package main

import (
	"context"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/client"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/server"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go client.DoAccrualIfPossible(logger, storage, cfg)

	router := server.NewRouter(logger, storage, cfg).Serve()

	server := &http.Server{Addr: cfg.GetHTTPServerAddress(), Handler: router}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-serverCtx.Done()
}
