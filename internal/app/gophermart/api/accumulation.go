package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/luhn"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strconv"
)

type AccumulationServiceWorker interface {
	AddOrder(accumulation *model.Accumulation) error
	GetAllOrders(userID int) ([]model.Accumulation, error)
	GetUserBalance(userID int) *model.Balance
	IsTheBalanceGreaterThanTheWriteOffAmount(userID int, amount float64) bool
	Withdraw(userID int, withdraw *model.Withdraw) error
	GetAllUserWithdrawals(userID int) *model.Withdrawals
	OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool
	OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool
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
			orderNumber  int
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
		orderNumber, err = strconv.Atoi(string(body))
		if err != nil {
			logger.Error("request format is incorrect", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "request format is incorrect")
			return
		}

		if !luhn.Valid(orderNumber) {
			logger.Error("order number format is incorrect")
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, "order number format is incorrect")
			return
		}

		accumulation.OrderNumber = orderNumber
		accumulation.UserID = r.Context().Value(model.ContextKeyUserID).(int)

		if a.accumulationService.OrderExistsAndBelongsToTheUser(&accumulation) {
			logger.Info("order exists and belongs to the user")
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, "order exists and belongs to the user")
			return
		}

		if a.accumulationService.OrderExistsAndDoesNotBelongToTheUser(&accumulation) {
			logger.Info("order exists and does not belong to the user")
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, "order exists and does not belong to the user")
			return
		}

		err = a.accumulationService.AddOrder(&accumulation)
		if err != nil {
			logger.Error("add user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "error when adding user")
			return
		}

		logger.Info("Put order:", accumulation)
		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, "order number accepted for processing")
	}
}

func (a *Accumulation) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetAllOrders"

		logger := a.logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value(model.ContextKeyUserID).(int)
		orders, err := a.accumulationService.GetAllOrders(userID)
		if err != nil {
			logger.Error("get all orders", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "error when getting all order numbers")
			return
		}

		if len(orders) == 0 {
			logger.Info("user have not order numbers")
			w.WriteHeader(http.StatusNoContent)
			render.JSON(w, r, "user have not order numbers")
			return
		}

		logger.Info("Order items:", orders)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
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

		w.Header().Add("content-type", "application/json")
		userID := r.Context().Value(model.ContextKeyUserID).(int)
		balance := a.accumulationService.GetUserBalance(userID)

		logger.Info("Current balance", balance)
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
		userID := r.Context().Value(model.ContextKeyUserID).(int)

		if !luhn.Valid(withdraw.OrderNumber) {
			logger.Error("order number format is incorrect")
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, "order number format is incorrect")
			return
		}

		if !a.accumulationService.IsTheBalanceGreaterThanTheWriteOffAmount(userID, withdraw.WriteOffAmount) {
			logger.Error("not enough points")
			w.WriteHeader(http.StatusPaymentRequired)
			render.JSON(w, r, "not enough points")
			return
		}

		err := a.accumulationService.Withdraw(userID, &withdraw)
		if err != nil {
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

		userID := r.Context().Value(model.ContextKeyUserID).(int)
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
