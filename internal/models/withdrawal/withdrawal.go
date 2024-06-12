package withdrawal

import "time"

type Withdrawal struct {
	ID          int
	Sum         float64
	ProcessedAt time.Time
	OrderNumber string
	BalanceID   int
}
