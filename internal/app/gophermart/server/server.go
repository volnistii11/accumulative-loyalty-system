package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api/accumulation"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api/auth"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"golang.org/x/exp/slog"
)

type Router struct {
	httpServer *chi.Mux
	logger     *slog.Logger
	db         *sqlx.DB
	cfg        *config.ParserGetter
}

func NewRouter(logger *slog.Logger, db *sqlx.DB, cfg *config.ParserGetter) *Router {
	return &Router{
		httpServer: chi.NewRouter(),
		logger:     logger,
		db:         db,
		cfg:        cfg,
	}
}

func (r *Router) Serve() *chi.Mux {
	apiAuth := auth.NewAuth()
	apiAccumulation := accumulation.NewAccumulation()

	r.httpServer.Use(middleware.RequestID)
	r.httpServer.Use(middleware.Logger)
	r.httpServer.Use(middleware.Recoverer)
	r.httpServer.Use(middleware.URLFormat)

	r.httpServer.Post("/api/user/register", apiAuth.RegisterUser())
	r.httpServer.Post("/api/user/login", apiAuth.AuthenticateUser())
	r.httpServer.Post("/api/user/orders", apiAccumulation.PutOrder())
	r.httpServer.Get("/api/user/orders", apiAccumulation.GetAllOrders())
	r.httpServer.Get("/api/user/balance", apiAccumulation.GetUserBalance())
	r.httpServer.Post("/api/user/balance/withdraw", apiAccumulation.DoWithdraw())
	r.httpServer.Get("/api/user/withdrawals", apiAccumulation.GetAllUserWithdrawls())

	return r.httpServer
}
