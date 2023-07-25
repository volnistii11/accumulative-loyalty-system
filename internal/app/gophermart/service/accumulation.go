package service

import (
	"github.com/volnistii11/accumulative-loyalty-system/internal/model"
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
	accumulation.UploadedAt = time.Now()

	err := db.AddOrder(accumulation)
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
