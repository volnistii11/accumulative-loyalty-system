package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/luhn"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strconv"
)

type Accumulation struct {
	accumulationService *service.Accumulation
}

func NewAccumulation(accumulationService *service.Accumulation) *Accumulation {
	return &Accumulation{
		accumulationService: accumulationService,
	}
}

func (a *Accumulation) PutOrder(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.PutOrder"
		var (
			err          error
			orderNumber  int
			accumulation model.Accumulation
		)

		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		defer r.Body.Close()
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

		if a.accumulationService.OrderExistsAndBelongsToTheUser(&accumulation, storage) {
			logger.Info("order exists and belongs to the user")
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, "order exists and belongs to the user")
			return
		}

		if a.accumulationService.OrderExistsAndDoesNotBelongToTheUser(&accumulation, storage) {
			logger.Info("order exists and does not belong to the user")
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, "order exists and does not belong to the user")
			return
		}

		err = a.accumulationService.AddOrder(&accumulation, storage)
		if err != nil {
			logger.Error("add user", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "error when adding user")
			return
		}

		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, "order number accepted for processing")
	}
}

func (a *Accumulation) GetAllOrders(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetAllOrders"

		logger := logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value(model.ContextKeyUserID).(int)
		orders, err := a.accumulationService.GetAllOrders(userID, storage)
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

func (a *Accumulation) GetUserBalance(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetUserBalance"

		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Add("content-type", "application/json")
		userID := r.Context().Value(model.ContextKeyUserID).(int)
		balance := a.accumulationService.GetUserBalance(userID, storage)

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, balance)
	}
}

func (a *Accumulation) DoWithdraw(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.DoWithdraw"

		logger = logger.With(
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

		if !a.accumulationService.IsTheBalanceGreaterThanTheWriteOffAmount(userID, withdraw.WriteOffAmount, storage) {
			logger.Error("not enough points")
			w.WriteHeader(http.StatusPaymentRequired)
			render.JSON(w, r, "not enough points")
			return
		}

		err := a.accumulationService.Withdraw(userID, &withdraw, storage)
		if err != nil {
			logger.Error("withdraw failed", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "withdraw failed")
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "")
	}
}

func (a *Accumulation) GetAllUserWithdrawals(logger *slog.Logger, storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const destination = "api.accumulation.GetAllUserWithdrawals"

		logger = logger.With(
			slog.String("destination", destination),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value(model.ContextKeyUserID).(int)
		withdrawals := a.accumulationService.GetAllUserWithdrawals(userID, storage)
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
