package service

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"math"
	"time"
)

type AdderGetterChecker interface {
	AddOrder(accumulation *model.Accumulation) error
	GetAllOrders(userID int) ([]model.Accumulation, error)
	GetUserBalance(userID int) *model.Balance
	Withdraw(accumulation *model.Accumulation) error
	GetAllUserWithdrawals(userID int) *model.Withdrawals
	OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool
	OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool
}

type Accumulation struct {
	db AdderGetterChecker
}

func NewAccumulation(db AdderGetterChecker) *Accumulation {
	return &Accumulation{
		db: db,
	}
}

func (accum *Accumulation) AddOrder(accumulation *model.Accumulation) error {
	currentTime := time.Now()
	accumulation.UploadedAt = &currentTime
	accumulation.ProcessingStatus = "NEW"

	err := accum.db.AddOrder(accumulation)
	if err != nil {
		return err
	}
	return nil
}

func (accum *Accumulation) GetAllOrders(userID int) ([]model.Accumulation, error) {
	orders, err := accum.db.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (accum *Accumulation) GetUserBalance(userID int) *model.Balance {
	balance := accum.db.GetUserBalance(userID)
	balance.Withdrawn = math.Abs(balance.Withdrawn)
	balance.Current = balance.Current - balance.Withdrawn
	return balance
}

func (accum *Accumulation) IsTheBalanceGreaterThanTheWriteOffAmount(userID int, amount float64) bool {
	balance := accum.db.GetUserBalance(userID)
	finalBalance := balance.Current + balance.Withdrawn
	return finalBalance >= amount
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

func (accum *Accumulation) OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool {
	return accum.db.OrderExistsAndBelongsToTheUser(accumulation)
}

func (accum *Accumulation) OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool {
	return accum.db.OrderExistsAndDoesNotBelongToTheUser(accumulation)
}
