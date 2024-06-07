package order

import "time"

type Order struct {
	ID         int
	Number     string    `json:"number"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	UserID     int
}

const (
	StatusNew        = "NEW"
	StatusRegistered = "REGISTERED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)
