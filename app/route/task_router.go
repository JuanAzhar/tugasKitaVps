package route

import (
	"tugaskita/features/task/handler"
	"tugaskita/features/task/repository"
	"tugaskita/features/task/service"
	userRepo "tugaskita/features/user/repository"
	userService "tugaskita/features/user/service"
	m "tugaskita/utils/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func TaskRouter(db *gorm.DB, e *echo.Group) {
	userRepository := userRepo.New(db)
	userUseCase := userService.New(userRepository)

	taskRepository := repository.NewTaskRepository(db, userRepository)
	taskUseCase := service.NewTaskService(taskRepository)
	taskController := handler.New(taskUseCase, userUseCase)

	user := e.Group("/user-task")
	user.GET("/:id", taskController.ReadSpecificTask, m.JWTMiddleware())
	user.GET("", taskController.ReadAllTask, m.JWTMiddleware())
	user.POST("", taskController.UploadTaskUser, m.JWTMiddleware())
	user.GET("/riwayat/:id", taskController.FindUserTaskById, m.JWTMiddleware())
	user.GET("/riwayat", taskController.ReadHistoryTaskUser, m.JWTMiddleware())

	user.POST("/request", taskController.UploadRequestTaskUser, m.JWTMiddleware())
	user.GET("/req-riwayat", taskController.FindAllRequestTaskHistory, m.JWTMiddleware())
	user.GET("/request/:id", taskController.FindUserTaskReqyId, m.JWTMiddleware())
	
	user.GET("/religion", taskController.FindAllReligionTaskUser, m.JWTMiddleware())
	user.GET("/religion/:id", taskController.ReadSpecificReligionTask, m.JWTMiddleware())
	user.GET("/religion/history", taskController.ReligionTaskHistoryUser, m.JWTMiddleware())
	user.POST("/religion", taskController.UploadTaskReligionUser, m.JWTMiddleware())

	user.POST("/religion-req", taskController.UploadReligionTaskRequest, m.JWTMiddleware())
	user.GET("/religion-req/history", taskController.FindAllReligionTaskRequestHistory, m.JWTMiddleware()) 
	user.GET("/religion-req/history/:id", taskController.FindSpesificReligionTaskRequest, m.JWTMiddleware()) 
	
	user.GET("/sum-clear", taskController.CountUserClearTask, m.JWTMiddleware())
	
	admin := e.Group("/admin-task")
	admin.GET("/:id", taskController.ReadSpecificTask, m.JWTMiddleware())
	admin.PUT("/:id", taskController.UpdateTask, m.JWTMiddleware())
	admin.GET("", taskController.ReadAllTask, m.JWTMiddleware())
	admin.POST("", taskController.AddTask, m.JWTMiddleware())
	admin.DELETE("/:id", taskController.DeleteTask, m.JWTMiddleware())

	admin.GET("/user", taskController.FindAllUserTask, m.JWTMiddleware())
	admin.GET("/user/request", taskController.FindAllUserRequestTask, m.JWTMiddleware())
	admin.PUT("/user/request/:id", taskController.UpdateTaskReqStatus, m.JWTMiddleware())
	admin.PUT("/user/:id", taskController.UpdateTaskStatus, m.JWTMiddleware())
	admin.GET("/user/:id", taskController.FindUserTaskById, m.JWTMiddleware())
	admin.GET("/user/request/:id", taskController.FindUserTaskReqyId, m.JWTMiddleware())

	admin.GET("/religion/:id", taskController.ReadSpecificReligionTask, m.JWTMiddleware())
	admin.PUT("/religion/:id", taskController.UpdateReligionTask, m.JWTMiddleware())
	admin.GET("/religion", taskController.ReadAllReligionTask, m.JWTMiddleware())
	admin.POST("/religion", taskController.AddReligionTask, m.JWTMiddleware())
	admin.DELETE("/religion/:id", taskController.DeleteReligionTask, m.JWTMiddleware())

	admin.GET("/religion/user", taskController.FindAllUserReligionTask, m.JWTMiddleware())
	admin.GET("/religion/user/:id", taskController.FindSpecificUserReligionTask, m.JWTMiddleware())
	admin.PUT("/religion/user/:id", taskController.UpdateReligionTaskStatus, m.JWTMiddleware())

	admin.GET("/religion/user-req", taskController.GetAllUserReligionTaskRequest, m.JWTMiddleware())
	admin.GET("/religion/user-req/:id", taskController.FindSpesificReligionTaskRequest, m.JWTMiddleware())
	admin.PUT("/religion/user-req/:id", taskController.UpdateTaskReligionReqStatus, m.JWTMiddleware())
}
