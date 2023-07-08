package main

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"log"
)

func main() {
	cfg := config.New()
	err := cfg.Parse()
	if err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	// TODO: init logger: slog

	// TODO: init storage: sqlx, pgx

	// TODO: init router: chi, render

	// TODO: run server
}
