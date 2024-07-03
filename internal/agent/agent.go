package agent

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	"lystem/internal/config"
	"lystem/internal/models/order"
	"lystem/internal/request"
	"lystem/internal/storage"
	"lystem/internal/usecase"
)

type Agent struct {
	storage           storage.Storage
	logger            *zap.SugaredLogger
	waitBeforePoll    time.Duration
	requestMaxRetries int
	url               string
}

func New(db storage.Storage, options config.Config, logger *zap.Logger) *Agent {
	return &Agent{
		storage:           db,
		logger:            logger.Sugar(),
		waitBeforePoll:    10 * time.Second,
		requestMaxRetries: options.RequestMaxRetries,
		url:               options.AccrualSystemAddress,
	}
}

func (p *Agent) StartOrdersPolling(ctx context.Context, wg *sync.WaitGroup) {
	ordersTimer := time.NewTimer(p.waitBeforePoll)

	for {
		select {
		case <-ordersTimer.C:
			p.PollOrdersInfo(ctx)
			ordersTimer.Reset(p.waitBeforePoll)
		case <-ctx.Done():
			p.logger.Info("4 Gracefully stop orders timer")
			ordersTimer.Stop()
			wg.Done()
			return
		}
	}
}

func (p *Agent) PollOrdersInfo(ctx context.Context) {
	orderUsecase := usecase.NewOrderUsecase(p.storage)
	orders, err := orderUsecase.FindAllAccrual(ctx)
	if err != nil {
		p.logger.Warn("failed to find orders", "error", err)
	}

	if len(orders) == 0 {
		p.logger.Info("no orders to request")
		// wait to collect to some orders to request
		time.Sleep(p.waitBeforePoll)
		return
	}

	for _, o := range orders {
		if ctx.Err() != nil {
			return
		}
		p.GetOneOrderInfo(ctx, &o, 0)
	}
}

func (p *Agent) GetOneOrderInfo(ctx context.Context, o *order.Order, retryCount int) {
	resp, err := http.Get(p.url + "/api/orders/" + o.Number)
	if err != nil {
		p.logger.Error("failed to get order", "error", err)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		p.saveOkOrder(ctx, o, resp.Body)
	case http.StatusNoContent:
		p.saveInvalidOrder(ctx, o)
	case http.StatusTooManyRequests:
		if retryCount >= p.requestMaxRetries {
			return
		}

		sleepTime := p.waitBeforePoll
		retryAfterStr := resp.Header.Get("Retry-After")
		if retryAfter, err := strconv.Atoi(retryAfterStr); err == nil {
			sleepTime = time.Duration(retryAfter) * time.Second
		}

		time.Sleep(sleepTime)
		p.GetOneOrderInfo(ctx, o, retryCount+1)
	case http.StatusInternalServerError:
		time.Sleep(p.waitBeforePoll)
		p.GetOneOrderInfo(ctx, o, retryCount+1)
	}

	if err = resp.Body.Close(); err != nil {
		p.logger.Error("failed to close response body", "error", err)
	}
}

func (p *Agent) saveOkOrder(ctx context.Context, o *order.Order, respBody io.ReadCloser) {
	ordersUsecase := usecase.NewOrderUsecase(p.storage)

	var orderInfo request.GetOrderRequest
	if err := json.NewDecoder(respBody).Decode(&orderInfo); err != nil {
		p.logger.Error("failed to decode order", "error", err)
		return
	}

	o.Status = orderInfo.Status
	if orderInfo.Accrual > 0 {
		o.Accrual = orderInfo.Accrual
	}

	if err := ordersUsecase.Update(ctx, o); err != nil {
		p.logger.Error("failed to update order", "error", err)
	}
}

func (p *Agent) saveInvalidOrder(ctx context.Context, o *order.Order) {
	ordersUsecase := usecase.NewOrderUsecase(p.storage)

	o.Status = order.StatusInvalid
	if err := ordersUsecase.Update(ctx, o); err != nil {
		p.logger.Error("failed to update order", "error", err)
	}
}
