package factory

import (
	"crypto/rand"
	"encoding/hex"

	"lystem/internal/models/user"
	"lystem/internal/request"
)

const (
	saltSize = 16
)

type UserFactory struct {
}

type UserFactoryMethods interface {
	Build(user request.CreateUser) (user.User, error)
}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (u *UserFactory) Build(userReq request.CreateUser) (*user.User, error) {
	var newUser user.User

	salt, err := generateRandomSalt()
	if err != nil {
		return nil, err
	}

	newUser.Login = userReq.Login
	newUser.SetHashedPassword(userReq.Password, salt)
	newUser.Salt = hex.EncodeToString(salt)
	return &newUser, nil
}

func generateRandomSalt() ([]byte, error) {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}
	return salt, nil
}
