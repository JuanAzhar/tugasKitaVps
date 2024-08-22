package dto

type PenaltyRequest struct {
	UserId      string `json:"user_id"`
	Point       int    `json:"point"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
