package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api/accumulation"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/api/auth"
)

type Router struct {
	httpServer *chi.Mux
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Serve() *chi.Mux {
	apiAuth := auth.NewAuth()
	apiAccumulation := accumulation.NewAccumulation()

	r.httpServer.Post("/api/user/register", apiAuth.RegisterUser())
	r.httpServer.Post("/api/user/login", apiAuth.AuthenticateUser())
	r.httpServer.Post("/api/user/orders", apiAccumulation.PutOrder())
	r.httpServer.Get("/api/user/orders", apiAccumulation.GetAllOrders())
	r.httpServer.Get("/api/user/balance", apiAccumulation.GetUserBalance())
	r.httpServer.Post("/api/user/balance/withdraw", apiAccumulation.DoWithdraw())
	r.httpServer.Get("/api/user/withdrawals", apiAccumulation.GetAllUserWithdrawls())

	return r.httpServer
}
