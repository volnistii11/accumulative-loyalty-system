package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/constants"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/luhn"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type AccumulationServiceWorker interface {
	AddOrder(w http.ResponseWriter, accumulation *model.Accumulation) (http.ResponseWriter, error)
	GetAllOrders(w http.ResponseWriter, userID int) (http.ResponseWriter, []model.Accumulation, error)
	GetUserBalance(userID int) *model.Balance
	Withdraw(userID int, withdraw *model.Withdraw) error
	GetAllUserWithdrawals(userID int) *model.Withdrawals
}

type Accumulation struct {
	accumulationService AccumulationServiceWorker
	logger              *slog.Logger
}

func NewAccumulation(accumulationService AccumulationServiceWorker, logger *slog.Logger) *Accumulation {
	return &Accumulation{
		accumulationService: accumulationService,
		logger:              logger,
	}
}

func (a *Accumulation) PutOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.PutOrder"
		var (
			err          error
			accumulation model.Accumulation
		)

		logger := a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}
		accumulation.OrderNumber = string(body)
		accumulation.UserID = getUserIDFromRequest(r)

		w, err = a.accumulationService.AddOrder(w, &accumulation)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}

		logger.Info("Put order:", accumulation)
		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, "order number accepted for processing")
	}
}

func (a *Accumulation) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromRequest(r)
		w, orders, err := a.accumulationService.GetAllOrders(w, userID)
		if err != nil {
			render.JSON(w, r, err.Error())
			return
		}

		render.JSON(w, r, orders)
	}
}

func (a *Accumulation) GetUserBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetUserBalance"

		logger := a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := getUserIDFromRequest(r)
		balance := a.accumulationService.GetUserBalance(userID)

		logger.Info("Current balance", balance)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, balance)
	}
}

func (a *Accumulation) DoWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.DoWithdraw"

		logger := a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var withdraw model.Withdraw
		if err := render.DecodeJSON(r.Body, &withdraw); err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}

		if !luhn.Valid(withdraw.OrderNumber) {
			logger.Error("order number format is incorrect")
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, "order number format is incorrect")
			return
		}

		userID := getUserIDFromRequest(r)
		err := a.accumulationService.Withdraw(userID, &withdraw)
		if err != nil {
			if errors.Is(err, cerrors.ErrDBNotEnoughCoins) {
				logger.Error("not enough points")
				w.WriteHeader(http.StatusPaymentRequired)
				render.JSON(w, r, "not enough points")
				return
			}
			logger.Error("withdraw failed", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "withdraw failed")
			return
		}
		logger.Info("Withdraw:", withdraw)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "")
	}
}

func (a *Accumulation) GetAllUserWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetAllUserWithdrawals"

		logger := a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := getUserIDFromRequest(r)
		withdrawals := a.accumulationService.GetAllUserWithdrawals(userID)
		if len(*withdrawals) == 0 {
			logger.Info("user have not withdrawals")
			w.WriteHeader(http.StatusNoContent)
			render.JSON(w, r, "user have not withdrawals")
			return
		}

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, withdrawals)
	}
}

func getUserIDFromRequest(r *http.Request) int {
	return r.Context().Value(constants.ContextKeyUserID).(int)
}
