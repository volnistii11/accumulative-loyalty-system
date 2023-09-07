package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func SetUpAccumulationService() *service.Accumulation {
	return service.NewAccumulation(nil)
}

func SetUpAPIAccumulation(accumService *service.Accumulation) *Accumulation {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return NewAccumulation(accumService, logger)
}

func TestDoWithdraw(t *testing.T) {
	doWithdraw := model.Withdraw{
		OrderNumber:    "123",
		WriteOffAmount: 751,
	}
	bodyIncorrectOrderNumber, _ := json.Marshal(doWithdraw)
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
			request: "/api/user/balance/withdraw",
			body:    []byte(``),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "incorrect order number",
			request: "/api/user/balance/withdraw",
			body:    bodyIncorrectOrderNumber,
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := SetUpRouter()
			accumService := SetUpAccumulationService()
			accumAPI := SetUpAPIAccumulation(accumService)

			router.Post(tt.request, accumAPI.DoWithdraw())
			req, _ := http.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(tt.body))
			req.Header.Add("content-type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
