package entity

import (
	"time"

	"github.com/google/uuid"
)

type PenaltyCore struct {
	Id          uuid.UUID `json:"id"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Point       int       `json:"point"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
