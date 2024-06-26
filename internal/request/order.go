package request

import (
	"errors"
)

type SaveOrderRequest struct {
	Number string
}

var (
	errInvalidOrderNumber = errors.New("неверный формат номера заказа")
)

func (s *SaveOrderRequest) Parse(body []byte) error {
	s.Number = string(body)
	return nil
}

func (s *SaveOrderRequest) Validate() error {
	if validLuhn(s.Number) {
		return nil
	}

	return errInvalidOrderNumber
}

type GetOrderRequest struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
