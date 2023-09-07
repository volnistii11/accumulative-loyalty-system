package api

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func SetUpRouter() *chi.Mux {
	return chi.NewRouter()
}

func SetUpAuthService() *service.Auth {
	return service.NewAuth(nil)
}

func SetUpAPIAuth(authService *service.Auth) *Auth {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return NewAuth(authService, logger)
}

func TestRegisterUser(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request string
		body    []byte
		want    want
	}{
		{
			name:    "body is empty",
			request: "/api/user/register",
			body:    []byte(``),
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := SetUpRouter()
			authService := SetUpAuthService()
			authAPI := SetUpAPIAuth(authService)

			router.Post(tt.request, authAPI.RegisterUser())
			req, _ := http.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(tt.body))
			req.Header.Add("content-type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		request string
		body    []byte
		want    want
	}{
		{
			name:    "body is empty",
			request: "/api/user/login",
			body:    []byte(``),
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := SetUpRouter()
			authService := SetUpAuthService()
			authAPI := SetUpAPIAuth(authService)

			router.Post(tt.request, authAPI.AuthenticateUser())
			req, _ := http.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(tt.body))
			req.Header.Add("content-type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
