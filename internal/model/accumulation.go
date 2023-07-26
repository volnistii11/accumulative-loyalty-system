package model

import "time"

type Accumulation struct {
	ID               int        `json:"id,omitempty" gorm:"primaryKey"`
	UserID           int        `json:"user_id,omitempty"`
	OrderNumber      int        `json:"order_number,omitempty"`
	UploadedAt       *time.Time `json:"uploaded_at,omitempty"`
	ProcessingStatus string     `json:"processing_status,omitempty"`
	AccrualStatus    string     `json:"accrual_status,omitempty"`
	Amount           float64    `json:"amount,omitempty"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
}

type Accumulations []*Accumulation

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	OrderNumber    int     `json:"order,string"`
	WriteOffAmount float64 `json:"sum"`
}

type Withdrawal struct {
	OrderNumber int        `json:"order"`
	Amount      float64    `json:"sum"`
	ProcessedAt *time.Time `json:"processed_at"`
}

type Withdrawals []*Withdrawal

type AccrualSystemAnswer struct {
	OrderNumber   int     `json:"order,string"`
	AccrualStatus string  `json:"status"`
	Amount        float64 `json:"accrual"`
}
