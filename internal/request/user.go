package request

import (
	"errors"
)

type CreateUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	//todo add password_confirmation
}

var (
	errNoLogin       = errors.New("логин должен быть более 3-х символов")
	errWrongPassword = errors.New("пароль должен быть более 8-ми символов")
)

func (cu CreateUser) Validate() error {
	if cu.Login == "" || len(cu.Login) < 3 {
		return errNoLogin
	}

	if cu.Password == "" || len(cu.Password) < 8 {
		return errWrongPassword
	}

	return nil
}
