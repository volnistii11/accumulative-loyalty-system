package service

import (
	"fmt"
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
	"math"
	"time"
)

// Accumulation TODO:
//
//	type Accumulation struct {
//		example OrderExample
//	}
//
//	type OrderExample interface {
//		OrderAdder
//		AllOrdersGetter
//		//....
//	}
//
//	func NewAccumulation(db OrderExample) *Accumulation {
//		return &Accumulation{
//			example: db,
//		}
//	}
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

	fmt.Println("Add order", accumulation)
	err := db.AddOrder(accumulation)
	if err != nil {
		return err
	}
	return nil
}

type AllOrdersGetter interface {
	GetAllOrders(userID int) (*model.Accumulations, error)
}

func (accum *Accumulation) GetAllOrders(userID int, db AllOrdersGetter) (*model.Accumulations, error) {
	orders, err := db.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}
	fmt.Println("GetAllORders", orders)
	return orders, nil
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
	return finalBalance >= amount
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
		ProcessedAt: &currentTime,
	}
	err := db.Withdraw(accumulation)
	if err != nil {
		return err
	}
	return nil
}

type AllUserWithdrawalsGetter interface {
	GetAllUserWithdrawals(userID int) *model.Withdrawals
}

func (accum *Accumulation) GetAllUserWithdrawals(userID int, db AllUserWithdrawalsGetter) *model.Withdrawals {
	withdrawals := db.GetAllUserWithdrawals(userID)
	return withdrawals
}

type OrderChecker interface {
	OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation) bool
	OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation) bool
}

func (accum *Accumulation) OrderExistsAndBelongsToTheUser(accumulation *model.Accumulation, db OrderChecker) bool {
	return db.OrderExistsAndBelongsToTheUser(accumulation)
}

func (accum *Accumulation) OrderExistsAndDoesNotBelongToTheUser(accumulation *model.Accumulation, db OrderChecker) bool {
	return db.OrderExistsAndDoesNotBelongToTheUser(accumulation)
}
