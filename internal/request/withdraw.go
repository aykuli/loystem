package request

import "fmt"

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (wr *WithdrawRequest) Validate() error {
	fmt.Println("validation", validLuhn(wr.Order))
	if validLuhn(wr.Order) {
		return nil
	}

	return errInvalidOrderNumber

}
