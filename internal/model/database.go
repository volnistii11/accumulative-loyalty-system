package model

type Withdrawal struct {
	OrderNumber int
	Amount      float64
	ProcessedAt int
}

type Withdrawals []*Withdrawals
