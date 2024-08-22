package dto

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Point       int    `json:"point"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Start_date  string `json:"start_date"`
	End_date    string `json:"end_date"`
}

type UserTaskUploadRequest struct {
	TaskId      string `json:"task_id" form:"task_id"`
	UserId      string `json:"user_id"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type UserReqTaskRequest struct {
	Title       string `json:"title" form:"title"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Point       int    `json:"point" form:"point"`
	Status      string `json:"status" form:"status"`
	Message     string `json:"message"`
}

type ReligionTaskRequest struct {
	Title       string `json:"title"`
	Religion    string `json:"religion"`
	Point       int    `json:"point"`
	Start_date  string `json:"start_date"`
	End_date    string `json:"end_date"`
	Description string `json:"description"`
}

type ReligionTaskUploadRequest struct{
	TaskId      string `json:"task_id" form:"task_id"`
	UserId      string `json:"user_id"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type UserReqReligionTaskRequest struct {
	Title       string `json:"title" form:"title"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Point       int    `json:"point" form:"point"`
	Status      string `json:"status" form:"status"`
	Message     string `json:"message"`
}