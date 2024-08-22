package model

import (
	"time"

	"github.com/google/uuid"
)

type Reward struct {
	ID        uuid.UUID `gorm:"type:varchar(50);primaryKey;not null" json:"id"`
	Name      string
	Stock     int
	Price     int
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRewardRequest struct {
	Id         uuid.UUID `gorm:"type:varchar(50);primaryKey;not null" json:"id"`
	RewardId   string
	Price      int
	Amount     int
	TotalPrice int
	UserId     string
	Status     string `gorm:"type:varchar(20);default:'Perlu Review'" json:"status"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
