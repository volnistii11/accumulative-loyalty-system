package model

import "time"

type Accumulation struct {
	ID               int        `json:"id,omitempty" gorm:"primaryKey"`
	UserID           int        `json:"user_id,omitempty"`
	OrderNumber      string     `json:"number,omitempty"`
	UploadedAt       *time.Time `json:"uploaded_at,omitempty"`
	ProcessingStatus string     `json:"status,omitempty"`
	AccrualStatus    string     `json:"accrual_status,omitempty"`
	Amount           float64    `json:"accrual,omitempty"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
}

type Accumulations []*Accumulation

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	OrderNumber    string  `json:"order"`
	WriteOffAmount float64 `json:"sum"`
}

type Withdrawal struct {
	OrderNumber string     `json:"order"`
	Amount      float64    `json:"sum"`
	ProcessedAt *time.Time `json:"processed_at"`
}

type Withdrawals []*Withdrawal

type AccrualSystemAnswer struct {
	OrderNumber   string  `json:"order"`
	AccrualStatus string  `json:"status"`
	Amount        float64 `json:"accrual"`
}
