package entity

import (
	"time"

	"github.com/google/uuid"
)

type RewardCore struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Price     int       `json:"price"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRewardRequestCore struct {
	Id         uuid.UUID `json:"id"`
	RewardId   string    `json:"reward_id"`
	RewardName string    `json:"reward_name"`
	Price      int       `json:"price"`
	TotalPrice int       `json:"total_price"`
	UserId     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	Status     string    `json:"status"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
