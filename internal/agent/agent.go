package agent

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"lystem/internal/models/order"
	"lystem/internal/request"
)

type GophermartAgent struct {
	url string
}

var ErrNoSuchOrder = errors.New("no such order")

func New(url string) *GophermartAgent {
	return &GophermartAgent{url}
}

func (a *GophermartAgent) GetOrderInfo(o *order.Order) (*order.Order, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.url, o.Number)
	req := fiber.Get(url)
	req.Set("Accept", "application/json")

	code, body, errs := req.Bytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	//todo if code == 429 reetry

	if code == fiber.StatusNoContent {
		//should we remove order if its not registred in system
		return nil, ErrNoSuchOrder
	}

	var orderInfo request.GetOrderRequest
	if err := json.Unmarshal(body, &orderInfo); err != nil {
		return nil, err
	}

	o.Status = orderInfo.Status
	o.Accrual = orderInfo.Accrual
	return o, nil
}
