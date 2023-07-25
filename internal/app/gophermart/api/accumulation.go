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
		accumulation.UserID = r.Context().Value("user_id").(int)

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

func (a *Accumulation) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func (a *Accumulation) GetUserBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func (a *Accumulation) DoWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func (a *Accumulation) GetAllUserWithdrawls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}