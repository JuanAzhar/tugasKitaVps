package dto

import "time"

type RewardResponse struct {
	Id        string
	Name      string
	Stock     int
	Price     int
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RewardRequestResponse struct {
	Id         string
	RewardId   string
	RewardName string
	UserId     string
	UserName   string
	Status     string
	Amount     string
	TotalPrice string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
