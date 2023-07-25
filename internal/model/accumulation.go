package model

import "time"

type Accumulation struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UserID           int       `json:"user_id""`
	OrderNumber      int       `json:"order_number"`
	UploadedAt       time.Time `json:"uploaded_at"`
	ProcessingStatus string    `json:"processing_status" gorm:"default:NEW"`
	AccrualStatus    string    `json:"accrual_status" gorm:"default:null"`
	Amount           int       `json:"amount" gorm:"default:null"`
	ProcessedAt      time.Time `json:"processed_at" gorm:"default:null"`
}

type Accumulations []*Accumulation

type Withdrawal struct {
	OrderNumber int
	Amount      float64
	ProcessedAt int
}

type Withdrawals []*Withdrawals
