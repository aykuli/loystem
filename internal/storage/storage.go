package storage

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
)

type Storage interface {
	CreateUser(ctx *fasthttp.RequestCtx, user *user.User) (*user.User, error)
	FindUserByLogin(ctx *fasthttp.RequestCtx, login string) (*user.User, error)
	CreateSession(ctx *fasthttp.RequestCtx, user *user.User) (*session.Session, error)
	DeleteSession(ctx *fasthttp.RequestCtx, user *user.User) error
}
