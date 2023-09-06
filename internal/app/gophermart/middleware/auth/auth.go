package auth

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/constants"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"golang.org/x/exp/slog"
	"net/http"
)

func ParseToken(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			const destination = "middleware.auth.ParseToken"
			logger := logger.With(
				slog.String("destination", destination),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			jwtToken, err := r.Cookie("jwtToken")
			if err != nil {
				logger.Info("jwt token is not found")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := service.GetUserID(jwtToken.Value)
			if err != nil {
				logger.Info("", sl.Err(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), constants.ContextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
