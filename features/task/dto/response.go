package dto

import "time"

type TaskResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Point       int    `json:"point"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Start_date  string `json:"startDate"`
	End_date    string `json:"endDate"`
	Description string `json:"description"`
}

type TaskResponseDetail struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Start_date  string    `json:"startDate"`
	End_date    string    `json:"endDate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type UserTaskUploadResponse struct {
	Id          string `json:"id"`
	TaskId      string `json:"task_id"`
	TaskName    string `json:"task_name"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type UserReqTaksResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Point       int    `json:"point"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type ReligionTaskResponse struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Religion string `json:"religion"`
	Point    int    `json:"point"`
}
