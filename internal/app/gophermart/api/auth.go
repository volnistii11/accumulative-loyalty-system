package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"net/http"
)

type UserAuthorize interface {
	RegisterUser(w http.ResponseWriter, user *model.User) (http.ResponseWriter, error)
	AuthenticateUser(w http.ResponseWriter, user *model.User) (http.ResponseWriter, error)
}

type Auth struct {
	authService UserAuthorize
	logger      *slog.Logger
}

func NewAuth(authService UserAuthorize, logger *slog.Logger) *Auth {
	return &Auth{
		authService: authService,
		logger:      logger,
	}
}

func (a *Auth) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.RegisterUser"

		a.logger = a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user model.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			a.logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}

		w, err = a.authService.RegisterUser(w, &user)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}

		w, err = a.authService.AuthenticateUser(w, &user)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}
		render.JSON(w, r, "you are registered")
	}
}

func (a *Auth) AuthenticateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.AuthenticateUser"
		a.logger = a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user model.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			a.logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}

		w, err = a.authService.AuthenticateUser(w, &user)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
