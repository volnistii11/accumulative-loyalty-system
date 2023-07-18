package model

type Accumulation struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	OrderNumber      int    `json:"order_number"`
	UploadedAt       int    `json:"uploaded_at"`
	ProcessingStatus string `json:"processing_status"`
	AccrualStatus    string `json:"accrual_status"`
	Amount           int    `json:"amount"`
	ProcessedAt      int    `json:"processed_at"`
}

type Accumulations []*Accumulation
