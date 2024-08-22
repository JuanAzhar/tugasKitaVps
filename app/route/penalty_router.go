package route

import (
	"tugaskita/features/penalty/handler"
	"tugaskita/features/penalty/repository"
	"tugaskita/features/penalty/service"
	userRepo "tugaskita/features/user/repository"
	userService "tugaskita/features/user/service"
	m "tugaskita/utils/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func PenaltyRouter(db *gorm.DB, e *echo.Group) {
	userRepository := userRepo.New(db)
	userUseCase := userService.New(userRepository)

	penaltyRepository := repository.NewPenaltyRepository(db)
	penaltyUseCase := service.NewPenaltyService(penaltyRepository, userRepository)
	penaltyController := handler.New(penaltyUseCase, userUseCase)
	
	admin := e.Group("/admin-penalty")
	admin.POST("", penaltyController.CreatePenalty, m.JWTMiddleware())
	admin.GET("", penaltyController.FindAllPenalty, m.JWTMiddleware())
	admin.GET("/:id", penaltyController.FindSpecificPenalty, m.JWTMiddleware())
	admin.PUT("/:id", penaltyController.UpdatePenalty, m.JWTMiddleware())
	admin.DELETE("/:id", penaltyController.DeletePenalty, m.JWTMiddleware())

	user := e.Group("/user-penalty")
	user.GET("/:id", penaltyController.FindSpecificPenalty, m.JWTMiddleware())
	user.GET("/history", penaltyController.FindAllPenaltyHistory, m.JWTMiddleware())

	e.GET("/sum-penalty", penaltyController.CountUserPenalty, m.JWTMiddleware())
}
