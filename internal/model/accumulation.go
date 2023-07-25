package model

import "time"

type Accumulation struct {
	ID               int        `json:"id,omitempty" gorm:"primaryKey"`
	UserID           int        `json:"user_id,omitempty"`
	OrderNumber      int        `json:"order_number,omitempty"`
	UploadedAt       *time.Time `json:"uploaded_at,omitempty"`
	ProcessingStatus string     `json:"processing_status,omitempty"`
	AccrualStatus    string     `json:"accrual_status,omitempty"`
	Amount           int        `json:"amount,omitempty"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
}

type Accumulations []*Accumulation

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	OrderNumber int
	Amount      float64
	ProcessedAt int
}

type Withdrawals []*Withdrawals
