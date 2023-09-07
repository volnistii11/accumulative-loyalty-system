package service

import (
	"errors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/luhn"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"golang.org/x/exp/slog"
	"math"
	"net/http"
	"time"
)

type AdderGetterChecker interface {
	AddOrder(accumulation *model.Accumulation) error
	GetAllOrders(userID int) ([]model.Accumulation, error)
	GetUserBalance(userID int) *model.Balance
	Withdraw(accumulation *model.Accumulation) error
	GetAllUserWithdrawals(userID int) *model.Withdrawals
}

type Accumulation struct {
	db     AdderGetterChecker
	logger *slog.Logger
}

func NewAccumulation(db AdderGetterChecker, logger *slog.Logger) *Accumulation {
	return &Accumulation{
		db:     db,
		logger: logger,
	}
}

func (accum *Accumulation) AddOrder(w http.ResponseWriter, accumulation *model.Accumulation) (http.ResponseWriter, error) {
	if !luhn.Valid(accumulation.OrderNumber) {
		accum.logger.Error("order number format is incorrect")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return w, cerrors.ErrOrderNumberIncorrect
	}

	currentTime := time.Now()
	accumulation.UploadedAt = &currentTime
	accumulation.ProcessingStatus = "NEW"

	err := accum.db.AddOrder(accumulation)
	if err != nil {
		if errors.Is(err, cerrors.ErrDBOrderExistsAndDoesNotBelongToTheUser) {
			accum.logger.Info("order exists and does not belong to the user")
			w.WriteHeader(http.StatusConflict)
			return w, cerrors.ErrDBOrderExistsAndDoesNotBelongToTheUser
		}

		if errors.Is(err, cerrors.ErrDBOrderExistsAndBelongsToTheUser) {
			accum.logger.Info("order exists and belongs to the user")
			w.WriteHeader(http.StatusOK)
			return w, cerrors.ErrDBOrderExistsAndBelongsToTheUser
		}

		accum.logger.Error("add user", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return w, err
	}
	return w, nil
}

func (accum *Accumulation) GetAllOrders(w http.ResponseWriter, userID int) (http.ResponseWriter, []model.Accumulation, error) {
	orders, err := accum.db.GetAllOrders(userID)
	if err != nil {
		accum.logger.Error("get all orders", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return w, nil, err
	}

	if len(orders) == 0 {
		accum.logger.Info("user have not order numbers")
		w.WriteHeader(http.StatusNoContent)
		return w, nil, cerrors.ErrHTTPStatusNoContent
	}

	accum.logger.Info("Order items:", orders)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	return w, orders, nil
}

func (accum *Accumulation) GetUserBalance(userID int) *model.Balance {
	balance := accum.db.GetUserBalance(userID)
	balance.Withdrawn = math.Abs(balance.Withdrawn)
	balance.Current = balance.Current - balance.Withdrawn
	return balance
}

func (accum *Accumulation) Withdraw(userID int, withdraw *model.Withdraw) error {
	currentTime := time.Now()
	accumulation := &model.Accumulation{
		UserID:      userID,
		OrderNumber: withdraw.OrderNumber,
		Amount:      -withdraw.WriteOffAmount,
		ProcessedAt: &currentTime,
	}
	err := accum.db.Withdraw(accumulation)
	if err != nil {
		return err
	}
	return nil
}

func (accum *Accumulation) GetAllUserWithdrawals(userID int) *model.Withdrawals {
	withdrawals := accum.db.GetAllUserWithdrawals(userID)
	return withdrawals
}
