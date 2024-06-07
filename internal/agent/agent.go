package agent

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"lystem/internal/config"
	"lystem/internal/models/order"
	"lystem/internal/request"
)

type GophermartAgent struct {
}

var ErrNoSuchOrder = errors.New("no such order")

func New() *GophermartAgent {
	return &GophermartAgent{}
}

func (a *GophermartAgent) GetOrderInfo(newOrder *order.Order) (*order.Order, error) {
	url := fmt.Sprintf("%s/api/orders/%s", config.Options.AccrualSystemAddress, newOrder.Number)
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

	newOrder.Status = orderInfo.Status
	newOrder.Accrual = orderInfo.Accrual
	return newOrder, nil
}

//
//func (a *Agent) Start() {
//	requestTicker := time.NewTicker(config.Options.PollInterval)
//	defer requestTicker.Stop()
//	for {
//		select {
//		case <-requestTicker.C:
//			err := makeReq(agent)
//			if err != nil {
//				log.Printf("agent error: %v", err)
//			}
//		}
//	}
//}
//
