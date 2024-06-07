package withdrawal

import "time"

type Withdrawal struct {
	ID          int
	Sum         float64
	ProcessedAt time.Time
	OrderID     int
	OrderNumber string
	BalanceID   int
}
