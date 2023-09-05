package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/gerr"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"net/http"
)

type UserAuthorize interface {
	RegisterUser(user *model.User) error
	AuthenticateUser(user *model.User) (string, error)
}

type Auth struct {
	authService UserAuthorize
}

func NewAuth(authService UserAuthorize) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (a *Auth) RegisterUser(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.RegisterUser"

		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user model.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}
		if user.Login == "" || user.Password == "" {
			logger.Error("wrong request format")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "wrong request format")
			return
		}

		err = a.authService.RegisterUser(&user)
		if err != nil {
			logger.Error("failed user register", sl.Err(err))
			if gerr.IsDuplicateKey(err) {
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, "user already exist")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal error")
			return
		}
		logger.Info("user registered")

		jwtToken, err := a.authService.AuthenticateUser(&user)
		if err != nil {
			logger.Error("failed user authentication", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal error")
			return
		}
		logger.Info("user authenticated")

		cookie := http.Cookie{Name: "jwtToken", Value: jwtToken}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "you are registered")
	}
}

func (a *Auth) AuthenticateUser(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.auth.AuthenticateUser"
		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user model.User
		err := render.DecodeJSON(r.Body, &user)
		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}
		if user.Login == "" || user.Password == "" {
			logger.Error("wrong request format")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "wrong request format")
			return
		}

		jwtToken, err := a.authService.AuthenticateUser(&user)
		if err != nil {
			logger.Error("failed user authentication", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal error")
			return
		}
		if jwtToken == "" {
			logger.Error("user or password is incorrect")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, "user or password is incorrect")
			return
		}
		logger.Info("user authenticated")

		cookie := http.Cookie{Name: "jwtToken", Value: jwtToken}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}
