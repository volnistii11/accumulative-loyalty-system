package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
)

type UserRegistrar interface {
	RegisterUser(login string, pass string) error
}

func RegisterUser(logger *slog.Logger, userRegister UserRegistrar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.RegisterUser"
		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		logger.Info("user registered")
	}
}

type UserAuthenticator interface {
	AuthenticateUser(login string, pass string) error
}

func AuthenticateUser(logger *slog.Logger, authenticator UserAuthenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.AuthenticateUser"
		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		logger.Info("user authenticated")
	}
}
