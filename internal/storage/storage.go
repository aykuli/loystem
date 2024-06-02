package storage

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
)

type Storage interface {
	CreateUser(ctx *fasthttp.RequestCtx, u *user.User) (*user.User, error)
	FindUserByLogin(ctx *fasthttp.RequestCtx, login string) (*user.User, error)
	FindUserByToken(ctx *fasthttp.RequestCtx, token string) (*user.User, error)
	CreateSession(ctx *fasthttp.RequestCtx, u *user.User) (*session.Session, error)
	DeleteSession(ctx *fasthttp.RequestCtx, u *user.User) error

	FindBalance(ctx *fasthttp.RequestCtx, u *user.User) (*balance.Balance, error)
	IncreaseBalance(ctx *fasthttp.RequestCtx, o *order.Order, u *user.User) error
	DeductFromBalance(ctx *fasthttp.RequestCtx, w *withdrawal.Withdrawal, u *user.User) error

	CreateWithdraw(ctx *fasthttp.RequestCtx, orderNumber string, u *user.User, sum float64) error
	FindWithdrawals(ctx *fasthttp.RequestCtx, balance *balance.Balance) ([]withdrawal.Withdrawal, error)

	FindOrderByNumber(ctx *fasthttp.RequestCtx, number string) (*order.Order, error)
	SaveOrder(ctx *fasthttp.RequestCtx, number string, userID int) (*order.Order, error)
	UpdateOrder(ctx *fasthttp.RequestCtx, o *order.Order) (*order.Order, error)
	SelectOrders(ctx *fasthttp.RequestCtx, u *user.User) ([]order.Order, error)
}
