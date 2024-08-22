package dto

type UserRequest struct {
	Name     string `json:"name" form:"name"`
	Image    string `json:"image" form:"image"`
	Address  string `json:"address" form:"address"`
	School   string `json:"school" form:"school"`
	Class    string `json:"class" form:"class"`
	Religion string `json:"religion" form:"religion"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Point    string `json:"point" form:"point"`
}
