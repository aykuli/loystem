package request

import "errors"

type CreateSession struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var errInvalidCreds = errors.New("неверный формат запроса")

func (s *CreateSession) Validate() error {
	if s.Login == "" || s.Password == "" {
		return errInvalidCreds
	}

	return nil
}
