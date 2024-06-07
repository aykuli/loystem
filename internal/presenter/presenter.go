package presenter

import (
	"time"

	"github.com/google/uuid"
)

type Common struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

type Session struct {
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
