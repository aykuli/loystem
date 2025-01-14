package presenter

import (
	"time"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/withdrawal"
)

type Common struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func NewSuccess(payload interface{}) *Common {
	return &Common{Success: true, Payload: payload}
}

func NewFailure(err error) *Common {
	return &Common{Success: false, Message: err.Error()}
}

type ResponseOrder struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrdersResponse(orders []order.Order) []ResponseOrder {
	var rOrders []ResponseOrder
	for _, o := range orders {
		rOrders = append(rOrders, ResponseOrder{Number: o.Number, Status: o.Status, Accrual: o.Accrual, UploadedAt: o.UploadedAt})
	}
	return rOrders
}

type ResponseBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceResponse(b *balance.Balance, ws []withdrawal.Withdrawal) ResponseBalance {
	var withdrawnSum float64
	for _, w := range ws {
		withdrawnSum += w.Sum
	}
	return ResponseBalance{Current: b.Current, Withdrawn: withdrawnSum}
}

type ResponseWithdrawals struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func NewWithdrawalsResponse(ws []withdrawal.Withdrawal) []ResponseWithdrawals {
	var responses []ResponseWithdrawals
	for _, w := range ws {
		responses = append(responses, ResponseWithdrawals{Order: w.OrderNumber, Sum: w.Sum, ProcessedAt: w.ProcessedAt})
	}
	return responses
}
