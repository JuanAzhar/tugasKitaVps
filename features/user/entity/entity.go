package entity

import "time"

type UserCore struct {
	ID         string    `json:"id"`
	Name       string    `json:"username"`
	Address    string    `json:"address"`
	School     string    `json:"school"`
	Class      string    `json:"class"`
	Image      string    `json:"image"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	Religion   string    `json:"religion"`
	Point      string    `json:"point"`
	TotalPoint string    `gorm:"Varchar(100);not null;default:user" json:"total_point"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"update_at"`
}

type UserPointCore struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Type      string    `json:"type"`
	TaskName  string    `json:"task_name"`
	Point     int       `json:"point"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}
