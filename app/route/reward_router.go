package route

import (
	"tugaskita/features/reward/handler"
	"tugaskita/features/reward/repository"
	"tugaskita/features/reward/service"

	userR "tugaskita/features/user/repository"
	userS "tugaskita/features/user/service"
	m "tugaskita/utils/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RewardRouter(db *gorm.DB, e *echo.Group) {
	userRepository := userR.New(db)
	userUseCase := userS.New(userRepository)

	rewardRepository := repository.NewRewardRepository(db, userRepository)
	rewardUseCase := service.NewRewardService(rewardRepository, userRepository)
	rewardController := handler.New(rewardUseCase, userUseCase)

	user := e.Group("/user-reward")
	user.GET("", rewardController.ReadAllReward, m.JWTMiddleware())
	user.GET("/:id", rewardController.ReadSpecificReward, m.JWTMiddleware())
	user.GET("/history", rewardController.FindAllRewardHistory, m.JWTMiddleware())
	user.POST("/exchange", rewardController.UploadRewardRequest, m.JWTMiddleware())

	admin := e.Group("/admin-reward")
	admin.GET("", rewardController.ReadAllReward, m.JWTMiddleware())
	admin.POST("", rewardController.AddReward, m.JWTMiddleware())
	admin.GET("/:id", rewardController.ReadSpecificReward, m.JWTMiddleware())
	admin.PUT("/:id", rewardController.UpdateReward, m.JWTMiddleware())
	admin.DELETE("/:id", rewardController.DeleteReward, m.JWTMiddleware())
	admin.GET("/user", rewardController.FindAllUploadReward, m.JWTMiddleware())
	admin.GET("/user/:id", rewardController.FindUserRewardById, m.JWTMiddleware())
	admin.PUT("/user/:id", rewardController.UpdateReqRewardStatus, m.JWTMiddleware())

}
