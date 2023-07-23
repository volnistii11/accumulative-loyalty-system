package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
)

type Router struct {
	httpServer *chi.Mux
	logger     *slog.Logger
	storage    *database.Storage
	cfg        *config.ParserGetter
}

func NewRouter(logger *slog.Logger, storage *database.Storage, cfg *config.ParserGetter) *Router {
	return &Router{
		httpServer: chi.NewRouter(),
		logger:     logger,
		storage:    storage,
		cfg:        cfg,
	}
}

func (r *Router) Serve() *chi.Mux {
	apiAccumulation := api.NewAccumulation()

	r.httpServer.Use(middleware.RequestID)
	r.httpServer.Use(middleware.Logger)
	r.httpServer.Use(middleware.Recoverer)
	r.httpServer.Use(middleware.URLFormat)

	r.httpServer.Post("/api/user/register", api.RegisterUser(r.logger, r.storage))
	r.httpServer.Post("/api/user/login", api.AuthenticateUser(r.logger, r.storage))
	r.httpServer.Post("/api/user/orders", apiAccumulation.PutOrder())
	r.httpServer.Get("/api/user/orders", apiAccumulation.GetAllOrders())
	r.httpServer.Get("/api/user/balance", apiAccumulation.GetUserBalance())
	r.httpServer.Post("/api/user/balance/withdraw", apiAccumulation.DoWithdraw())
	r.httpServer.Get("/api/user/withdrawals", apiAccumulation.GetAllUserWithdrawls())

	return r.httpServer
}
