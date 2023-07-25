package service

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"math"
	"time"
)

type Accumulation struct {
}

func NewAccumulation() *Accumulation {
	return &Accumulation{}
}

type OrderAdder interface {
	AddOrder(accumulation *model.Accumulation) error
}

func (accum *Accumulation) AddOrder(accumulation *model.Accumulation, db OrderAdder) error {
	currentTime := time.Now()
	accumulation.UploadedAt = &currentTime
	accumulation.ProcessingStatus = "NEW"

	err := db.AddOrder(accumulation)
	if err != nil {
		return err
	}
	return nil
}

type AllOrdersGetter interface {
	GetAllOrders(userID int) *model.Accumulations
}

func (accum *Accumulation) GetAllOrders(userID int, db AllOrdersGetter) *model.Accumulations {
	orders := db.GetAllOrders(userID)
	return orders
}

type UserBalanceGetter interface {
	GetUserBalance(userID int) *model.Balance
}

func (accum *Accumulation) GetUserBalance(userID int, db UserBalanceGetter) *model.Balance {
	balance := db.GetUserBalance(userID)
	balance.Withdrawn = math.Abs(balance.Withdrawn)
	return balance
}

func (accum *Accumulation) IsTheBalanceGreaterThanTheWriteOffAmount(userID int, amount float64, db UserBalanceGetter) bool {
	balance := db.GetUserBalance(userID)
	finalBalance := balance.Current + balance.Withdrawn
	if finalBalance >= amount {
		return true
	}
	return false
}

type PointsWithdrawal interface {
	Withdraw(accumulation *model.Accumulation) error
}

func (accum *Accumulation) Withdraw(userID int, withdraw *model.Withdraw, db PointsWithdrawal) error {
	currentTime := time.Now()
	accumulation := &model.Accumulation{
		UserID:      userID,
		OrderNumber: withdraw.OrderNumber,
		Amount:      -withdraw.WriteOffAmount,
		UploadedAt:  &currentTime,
	}
	err := db.Withdraw(accumulation)
	if err != nil {
		return err
	}
	return nil
}

type OrderChecker interface {
	OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool
	OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool
}

func (accum *Accumulation) OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation, db OrderChecker) bool {
	if db.OrderExistsAndBelongsToTheUser(accumulation) {
		return true
	}
	return false
}

func (accum *Accumulation) OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation, db OrderChecker) bool {
	if db.OrderExistsAndDoesNotBelongToTheUser(accumulation) {
		return true
	}
	return false
}
