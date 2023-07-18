package accumulation

import (
	"net/http"
)

type Accumulation struct {
}

func NewAccumulation() *Accumulation {
	return &Accumulation{}
}

func (a *Accumulation) PutOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: put in order into database with postgres.PutOrder
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
