package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskCore struct {
	ID          uuid.UUID `json:"id"`
	AdminId     string    `json:"admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Start_date  string    `json:"start_date"`
	End_date    string    `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserTaskUploadCore struct {
	Id          uuid.UUID `json:"id"`
	TaskId      string    `json:"task_id"`
	TaskName    string    `json:"task_name"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserTaskSubmissionCore struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Image       string    `json:"image"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReligionTaskCore struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Religion    string    `json:"religion"`
	Point       int       `json:"point"`
	Type        string    `json:"type"`
	Start_date  string    `json:"start_date"`
	End_date    string    `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserReligionTaskUploadCore struct {
	Id          uuid.UUID `json:"id"`
	TaskId      string    `json:"task_id"`
	TaskName    string    `json:"task_name"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"username"`
	Image       string    `json:"image"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserReligionReqTaskCore struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Image       string    `json:"image"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}