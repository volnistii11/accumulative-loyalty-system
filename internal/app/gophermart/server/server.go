package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/middleware/auth"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
)

type Router struct {
	httpServer *chi.Mux
	logger     *slog.Logger
	storage    *database.Storage
	cfg        config.ParserGetter
}

func NewRouter(logger *slog.Logger, storage *database.Storage, cfg config.ParserGetter) *Router {
	return &Router{
		httpServer: chi.NewRouter(),
		logger:     logger,
		storage:    storage,
		cfg:        cfg,
	}
}

func (router *Router) Serve() *chi.Mux {
	authService := service.NewAuth(router.storage, router.logger)
	apiAuth := api.NewAuth(authService, router.logger)

	accumulationService := service.NewAccumulation(router.storage, router.logger)
	apiAccumulation := api.NewAccumulation(accumulationService, router.logger)

	router.httpServer.Group(func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)

		r.Post("/api/user/register", apiAuth.RegisterUser())
		r.Post("/api/user/login", apiAuth.AuthenticateUser())
		r.Group(func(r chi.Router) {
			r.Use(auth.ParseToken(router.logger))
			r.Post("/api/user/orders", apiAccumulation.PutOrder())
			r.Get("/api/user/orders", apiAccumulation.GetAllOrders())
			r.Get("/api/user/balance", apiAccumulation.GetUserBalance())
			r.Post("/api/user/balance/withdraw", apiAccumulation.DoWithdraw())
			r.Get("/api/user/withdrawals", apiAccumulation.GetAllUserWithdrawals())
		})
	})

	return router.httpServer
}
