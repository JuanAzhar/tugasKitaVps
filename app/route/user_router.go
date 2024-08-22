package route

import (
	"tugaskita/features/user/handler"
	"tugaskita/features/user/repository"
	"tugaskita/features/user/service"
	m "tugaskita/utils/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRouter(db *gorm.DB, e *echo.Group) {
	userRepository := repository.New(db)
	userUseCase := service.New(userRepository)
	userController := handler.New(userUseCase)

	e.POST("/register", userController.Register)
	e.POST("/login", userController.Login)
	e.GET("", userController.ReadAllUser, m.JWTMiddleware())
	e.GET("/profile", userController.ReadProfileUser, m.JWTMiddleware())
	e.GET("/:id", userController.ReadSpecificUser, m.JWTMiddleware())
	e.DELETE("/:id", userController.DeleteUser, m.JWTMiddleware())
	e.GET("/rank", userController.GetRankUser, m.JWTMiddleware())
	e.PUT("/change-password", userController.ChangePassword, m.JWTMiddleware())
	e.PUT("/:id", userController.UpdateSiswa, m.JWTMiddleware())

	e.POST("/monthly-reset", userController.MonthlyResetPoint, m.JWTMiddleware())
	e.POST("/annual-reset", userController.AnnualResetPoint, m.JWTMiddleware())

	e.GET("/user-point-history", userController.GetAllUserPointHistory, m.JWTMiddleware())
	e.GET("/user-point-history/:id", userController.GetSpecificUserPointHistory, m.JWTMiddleware())
	e.GET("/point-history",userController.GetUserPointHistory, m.JWTMiddleware())
}
