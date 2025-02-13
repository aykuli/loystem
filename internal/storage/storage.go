package storage

import (
	"context"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/session"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
)

type Storage interface {
	CreateUser(ctx context.Context, u *user.User) (*user.User, error)
	FindUserByLogin(ctx context.Context, login string) (*user.User, error)
	FindUserByToken(ctx context.Context, token string) (*user.User, error)
	CreateSession(ctx context.Context, u *user.User) (*session.Session, error)
	DeleteSession(ctx context.Context, u *user.User) error

	FindBalance(ctx context.Context, u *user.User) (*balance.Balance, error)
	DeductFromBalance(ctx context.Context, w *withdrawal.Withdrawal, u *user.User) error

	CreateWithdrawal(ctx context.Context, orderNumber string, u *user.User, sum float64) (*withdrawal.Withdrawal, error)
	FindWithdrawals(ctx context.Context, balance *balance.Balance) ([]withdrawal.Withdrawal, error)

	FindOrderByNumber(ctx context.Context, number string) (*order.Order, error)
	SaveOrder(ctx context.Context, number string, userID int) (*order.Order, error)
	UpdateOrder(ctx context.Context, o *order.Order) error
	UpdateOrderAndIncreaseBalance(ctx context.Context, o *order.Order) error
	FindAllUserOrders(ctx context.Context, u *user.User) ([]order.Order, error)
	SelectUnprocessedOrders(ctx context.Context, limit int) ([]order.Order, error)
}
