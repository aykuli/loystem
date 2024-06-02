package order

import "time"

type Order struct {
	ID         int
	Number     string
	Accrual    float64
	UploadedAt time.Time
	Status     string
	UserID     int
}

const (
	StatusNew        = "NEW"
	StatusRegistered = "REGISTERED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)
