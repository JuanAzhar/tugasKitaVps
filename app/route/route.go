package route

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(e *echo.Echo, db *gorm.DB) {
	user := e.Group("user")
	base := e.Group("")

	UserRouter(db, user)
	TaskRouter(db, base)
	RewardRouter(db, base)
	PenaltyRouter(db, base)
}
