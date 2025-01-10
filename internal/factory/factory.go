package factory

import (
	"lystem/internal/models/user"
	"lystem/internal/request"
)

type UserFactory struct {
	salt string
}

func NewUserFactory(salt string) *UserFactory {
	return &UserFactory{salt}
}

func (u *UserFactory) Build(userReq request.CreateUser) (*user.User, error) {
	var newUser user.User

	newUser.Login = userReq.Login
	newUser.SetHashedPassword(userReq.Password, u.salt)
	return &newUser, nil
}
