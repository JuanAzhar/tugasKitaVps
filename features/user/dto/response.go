package dto

type UserResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	School     string `json:"school"`
	Class      string `json:"class"`
	Image      string `json:"image"`
	Role       string `json:"role"`
	Religion   string `json:"religion"`
	Email      string `json:"email"`
	Point      string `json:"point"`
	TotalPoint string `json:"total_point"`
}

type UserRankResponse struct {
	Name  string `json:"name"`
	Point string `json:"point"`
}
