package request

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (wr *WithdrawRequest) Validate() error {
	return nil
	//if validLuhn(wr.Order) {
	//	return nil
	//}
	//
	//return errInvalidOrderNumber

}
