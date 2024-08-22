package handler

import (
	"net/http"
	"tugaskita/features/task/dto"
	"tugaskita/features/task/entity"
	user "tugaskita/features/user/entity"
	middleware "tugaskita/utils/jwt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskController struct {
	taskUsecase entity.TaskUseCaseInterface
	userUsecase user.UserUseCaseInterface
}

func New(taskUC entity.TaskUseCaseInterface, userUC user.UserUseCaseInterface) *TaskController {
	return &TaskController{
		taskUsecase: taskUC,
		userUsecase: userUC,
	}
}

func (handler *TaskController) AddTask(e echo.Context) error {
	userId, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	input := new(dto.TaskRequest)
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	data := entity.TaskCore{
		AdminId:     userId,
		Title:       input.Title,
		Description: input.Description,
		Point:       input.Point,
		Start_date:  input.Start_date,
		End_date:    input.End_date,
	}

	errTask := handler.taskUsecase.CreateTask(data)
	if errTask != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error create task",
			"error":   errTask.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes create task",
	})
}

func (handler *TaskController) ReadAllTask(e echo.Context) error {
	userId, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role == "admin" {
		data, err := handler.taskUsecase.FindAllTask()
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]any{
				"message": "error get all task",
			})
		}

		dataList := []dto.TaskResponse{}
		for _, v := range data {
			result := dto.TaskResponse{
				Id:          v.ID.String(),
				Title:       v.Title,
				Point:       v.Point,
				Status:      v.Status,
				Type:        v.Type,
				Start_date:  v.Start_date,
				End_date:    v.End_date,
				Description: v.Description,
			}
			dataList = append(dataList, result)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message": "get all admin task",
			"data":    dataList,
		})

	} else if role == "user" {
		data, err := handler.taskUsecase.FindTasksNotClaimedByUser(userId)
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]any{
				"message": "error get all task",
			})
		}

		dataList := []dto.TaskResponse{}
		for _, v := range data {
			result := dto.TaskResponse{
				Id:          v.ID.String(),
				Title:       v.Title,
				Point:       v.Point,
				Status:      v.Status,
				Type:        v.Type,
				Start_date:  v.Start_date,
				End_date:    v.End_date,
				Description: v.Description,
			}
			dataList = append(dataList, result)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message": "get all user task",
			"data":    dataList,
		})
	}

	return e.JSON(http.StatusBadRequest, map[string]any{
		"message": "access denied",
	})
}

func (handler *TaskController) ReadSpecificTask(e echo.Context) error {

	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "task not found",
		})
	}

	data, err := handler.taskUsecase.FindById(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific task",
		})
	}

	response := dto.TaskResponseDetail{
		Title:       data.Title,
		Description: data.Description,
		Point:       data.Point,
		Type:        data.Type,
		Message:     data.Message,
		Status:      data.Status,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get task",
		"data":    response,
	})
}

func (handler *TaskController) DeleteTask(e echo.Context) error {
	_, role, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")
	err := handler.taskUsecase.DeleteTask(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error deleting task",
			"error":   err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task deleted successfully",
	})
}

func (handler *TaskController) UpdateTask(e echo.Context) error {
	adminId, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := new(dto.TaskRequest)
	if errBind := e.Bind(data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	taskData := entity.TaskCore{
		AdminId:     adminId,
		Title:       data.Title,
		Description: data.Description,
		Point:       data.Point,
		Status:      data.Status,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
	}

	errUpdate := handler.taskUsecase.UpdateTask(idParams, taskData)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating task",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "task updated successfully",
	})
}

func (handler *TaskController) UpdateTaskStatus(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := dto.UserTaskUploadRequest{}
	if errBind := e.Bind(&data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	//get userId & taskId
	dataTask, _ := handler.taskUsecase.FindUserTaskById(idParams)

	status := entity.UserTaskUploadCore{
		TaskId:      dataTask.TaskId,
		UserId:      dataTask.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
	}

	errUpdate := handler.taskUsecase.UpdateTaskStatus(idParams, status)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating task status",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "task status updated",
	})
}

func (handler *TaskController) UpdateTaskReqStatus(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := dto.UserReqTaskRequest{}
	if errBind := e.Bind(&data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	//get userId & taskId
	dataTask, _ := handler.taskUsecase.FindUserTaskReqById(idParams)

	status := entity.UserTaskSubmissionCore{
		UserId:      dataTask.UserId,
		UserName:    dataTask.UserName,
		Title:       data.Title,
		Image:       data.Image,
		Description: data.Description,
		Point:       data.Point,
		Status:      data.Status,
		Message:     data.Message,
	}

	errUpdate := handler.taskUsecase.UpdateTaskReqStatus(idParams, status)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating request task status",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "task request status updated",
	})
}

func (handler *TaskController) ReadHistoryTaskUser(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.taskUsecase.FindAllClaimedTask(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get history task",
		})
	}

	dataList := []entity.UserTaskUploadCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)
		taskData, _ := handler.taskUsecase.FindById(v.TaskId)

		result := entity.UserTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			TaskName:    taskData.Title,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Type:        taskData.Type,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all task history",
		"data":    dataList,
	})
}

func (handler *TaskController) UploadTaskUser(e echo.Context) error {
	input := dto.UserTaskUploadRequest{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	image, err := e.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "No file uploaded",
			})
		}
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error uploading file",
		})
	}

	userId, _, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	dataInput := entity.UserTaskUploadCore{
		TaskId:      input.TaskId,
		Image:       input.Image,
		Description: input.Description,
	}
	dataInput.UserId = userId

	errUpload := handler.taskUsecase.UploadTask(dataInput, image)
	if errUpload != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error upload task",
			"error":   errUpload.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes upload task",
	})
}

func (handler *TaskController) FindAllUserTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	data, err := handler.taskUsecase.FindAllUserTask()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user task",
		})
	}

	dataList := []entity.UserTaskUploadCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)
		taskData, _ := handler.taskUsecase.FindById(v.TaskId)

		result := entity.UserTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			TaskName:    taskData.Title,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Type:        v.Type,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user task",
		"data":    dataList,
	})
}

func (handler *TaskController) FindUserTaskById(e echo.Context) error {
	idParams := e.Param("id")

	data, err := handler.taskUsecase.FindUserTaskById(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific task",
		})
	}

	userData, _ := handler.userUsecase.ReadSpecificUser(data.UserId)
	taskData, _ := handler.taskUsecase.FindById(data.TaskId)

	response := entity.UserTaskUploadCore{
		Id:          data.Id,
		UserId:      data.UserId,
		UserName:    userData.Name,
		TaskId:      data.TaskId,
		TaskName:    taskData.Title,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Type:        data.Type,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get task",
		"data":    response,
	})
}

func (handler *TaskController) FindUserTaskReqyId(e echo.Context) error {
	idParams := e.Param("id")

	data, err := handler.taskUsecase.FindUserTaskReqById(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific task",
		})
	}

	userData, _ := handler.userUsecase.ReadSpecificUser(data.UserId)

	response := entity.UserTaskSubmissionCore{
		Id:          data.Id,
		Title:       data.Title,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Point:       data.Point,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Type:        data.Type,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get task",
		"data":    response,
	})
}

func (handler *TaskController) UploadRequestTaskUser(e echo.Context) error {
	input := dto.UserReqTaskRequest{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	userId, _, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	image, err := e.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "No file uploaded",
			})
		}
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error uploading file",
		})
	}

	dataInput := entity.UserTaskSubmissionCore{
		Title:       input.Title,
		Point:       input.Point,
		Image:       input.Image,
		Description: input.Description,
	}
	dataInput.UserId = userId

	errUpload := handler.taskUsecase.UploadTaskRequest(dataInput, image)
	if errUpload != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error upload request task",
			"error":   errUpload.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes upload request task",
	})
}

func (handler *TaskController) FindAllUserRequestTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	data, err := handler.taskUsecase.FindAllRequestTask()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user request task",
		})
	}

	dataList := []entity.UserTaskSubmissionCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)

		result := entity.UserTaskSubmissionCore{
			Id:          v.Id,
			Title:       v.Title,
			Point:       v.Point,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Type:        v.Type,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user task request",
		"data":    dataList,
	})
}

func (handler *TaskController) FindAllRequestTaskHistory(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.taskUsecase.FindAllRequestTaskHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user request task history",
		})
	}

	dataList := []entity.UserTaskSubmissionCore{}
	for _, v := range data {
		result := entity.UserTaskSubmissionCore{
			Id:          v.Id,
			Title:       v.Title,
			Point:       v.Point,
			UserId:      v.UserId,
			Image:       v.Image,
			Description: v.Description,
			Message:     v.Message,
			Status:      v.Status,
			Type:        v.Type,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user task request history",
		"data":    dataList,
	})
}

func (handler *TaskController) CountUserClearTask(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	count, err := handler.taskUsecase.CountUserClearTask(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get cleared task sum",
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all task cleared sum",
		"count":   count,
	})
}

func (handler *TaskController) AddReligionTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	input := new(dto.ReligionTaskRequest)
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	data := entity.ReligionTaskCore{
		Title:       input.Title,
		Point:       input.Point,
		Religion:    input.Religion,
		Start_date:  input.Start_date,
		End_date:    input.End_date,
		Description: input.Description,
	}

	errTask := handler.taskUsecase.CreateTaskReligion(data)
	if errTask != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error create religion task",
			"error":   errTask.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes create religion task",
	})
}

func (handler *TaskController) ReadAllReligionTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role == "admin" {
		data, err := handler.taskUsecase.FindAllTaskReligion()
		if err != nil {
			return e.JSON(http.StatusBadRequest, map[string]any{
				"message": "error get all task",
			})
		}

		dataList := []entity.ReligionTaskCore{}
		for _, v := range data {
			result := entity.ReligionTaskCore{
				Id:          v.Id,
				Title:       v.Title,
				Description: v.Description,
				Type:        v.Type,
				Start_date:  v.Start_date,
				End_date:    v.End_date,
				Point:       v.Point,
				Religion:    v.Religion,
			}
			dataList = append(dataList, result)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message": "get all admin task",
			"data":    dataList,
		})

	} else if role == "user" {
		// data, err := handler.taskUsecase.FindTasksNotClaimedByUser(userId)
		// if err != nil {
		// 	return e.JSON(http.StatusBadRequest, map[string]any{
		// 		"message": "error get all task",
		// 	})
		// }

		// dataList := []dto.TaskResponse{}
		// for _, v := range data {
		// 	result := dto.TaskResponse{
		// 		Id:         v.ID.String(),
		// 		Title:      v.Title,
		// 		Point:      v.Point,
		// 		Status:     v.Status,
		// 		Type:       v.Type,
		// 		Start_date: v.Start_date,
		// 		End_date:   v.End_date,
		// 	}
		// 	dataList = append(dataList, result)
		// }

		return e.JSON(http.StatusOK, map[string]any{
			"message": "get all user task",
			// "data":    dataList,
		})
	}

	return e.JSON(http.StatusBadRequest, map[string]any{
		"message": "access denied",
	})
}

func (handler *TaskController) ReadSpecificReligionTask(e echo.Context) error {

	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "task not found",
		})
	}

	data, err := handler.taskUsecase.FindByIdReligionTask(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific task",
		})
	}

	response := entity.ReligionTaskCore{
		Id:          data.Id,
		Title:       data.Title,
		Description: data.Description,
		Type:        data.Type,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
		Point:       data.Point,
		Religion:    data.Religion,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get task",
		"data":    response,
	})
}

func (handler *TaskController) DeleteReligionTask(e echo.Context) error {
	_, role, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")
	err := handler.taskUsecase.DeleteTaskReligion(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error deleting task",
			"error":   err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task deleted successfully",
	})
}

func (handler *TaskController) UpdateReligionTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := new(dto.ReligionTaskRequest)
	if errBind := e.Bind(data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	taskData := entity.ReligionTaskCore{
		Title:    data.Title,
		Religion: data.Religion,
		Point:    data.Point,
	}

	errUpdate := handler.taskUsecase.UpdateTaskReligion(idParams, taskData)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating task",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "task updated successfully",
	})
}

func (handler *TaskController) FindAllReligionTaskUser(e echo.Context) error {
	userId, _, religion, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.taskUsecase.FindAllReligionTaskUser(religion, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all religion task",
		})
	}

	dataList := []entity.ReligionTaskCore{}
	for _, v := range data {
		result := entity.ReligionTaskCore{
			Id:          v.Id,
			Title:       v.Title,
			Point:       v.Point,
			Description: v.Description,
			Religion:    v.Description,
			Start_date:  v.Start_date,
			End_date:    v.End_date,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user religion task",
		"data":    dataList,
	})
}

func (handler *TaskController) ReligionTaskHistoryUser(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.taskUsecase.FindAllReligionTaskHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get religion task history",
		})
	}

	dataList := []entity.UserReligionTaskUploadCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)
		taskData, _ := handler.taskUsecase.FindByIdReligionTask(v.TaskId)

		result := entity.UserReligionTaskUploadCore{
			Id:          v.Id,
			UserId:      v.UserId,
			UserName:    userData.Name,
			TaskId:      v.TaskId,
			TaskName:    taskData.Title,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user religion task history",
		"data":    dataList,
	})
}

func (handler *TaskController) UploadTaskReligionUser(e echo.Context) error {
	input := dto.ReligionTaskUploadRequest{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	userId, _, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	image, err := e.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "No file uploaded",
			})
		}
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error uploading file",
		})
	}

	dataInput := entity.UserReligionTaskUploadCore{
		TaskId:      input.TaskId,
		Description: input.Description,
		Image:       input.Image,
	}
	dataInput.UserId = userId

	errUpload := handler.taskUsecase.UploadTaskReligion(dataInput, image)
	if errUpload != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error upload religion task",
			"error":   errUpload.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes upload religion task",
	})
}

func (handler *TaskController) FindAllUserReligionTask(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	data, err := handler.taskUsecase.FindAllUserReligionTaskUpload()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user religion task",
		})
	}

	dataList := []entity.UserReligionTaskUploadCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)
		taskData, _ := handler.taskUsecase.FindByIdReligionTask(v.TaskId)

		result := entity.UserReligionTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			TaskName:    taskData.Title,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Type:        v.Type,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user religion task",
		"data":    dataList,
	})
}

func (handler *TaskController) FindSpecificUserReligionTask(e echo.Context) error {
	idParams := e.Param("id")

	data, err := handler.taskUsecase.FindSpecificUserReligionTaskUpload(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific religion task",
		})
	}

	userData, _ := handler.userUsecase.ReadSpecificUser(data.UserId)
	taskData, _ := handler.taskUsecase.FindByIdReligionTask(data.TaskId)

	response := entity.UserReligionTaskUploadCore{
		Id:          data.Id,
		UserId:      data.UserId,
		UserName:    userData.Name,
		TaskId:      data.TaskId,
		TaskName:    taskData.Title,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Type:        data.Type,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get specific religion task",
		"data":    response,
	})
}

func (handler *TaskController) UpdateReligionTaskStatus(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := dto.ReligionTaskUploadRequest{}
	if errBind := e.Bind(&data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	//get userId & taskId
	dataTask, _ := handler.taskUsecase.FindSpecificUserReligionTaskUpload(idParams)

	status := entity.UserReligionTaskUploadCore{
		Id:          dataTask.Id,
		UserId:      dataTask.UserId,
		TaskId:      data.TaskId,
		TaskName:    dataTask.TaskName,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
	}

	errUpdate := handler.taskUsecase.UpdateReligionTaskStatus(dataTask.TaskId, status)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating religion task status",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "religion task status updated",
	})
}

func (handler *TaskController) FindAllReligionTaskRequestHistory(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.taskUsecase.FindAllReligionTaskRequestHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get religion task request history",
		})
	}

	dataList := []entity.UserReligionReqTaskCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)

		result := entity.UserReligionReqTaskCore{
			Id:          v.Id,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Title:       v.Title,
			Type:        v.Type,
			Point:       v.Point,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user task religion request",
		"data":    dataList,
	})
}

func (handler *TaskController) FindSpesificReligionTaskRequest(e echo.Context) error {
	idParams := e.Param("id")

	data, err := handler.taskUsecase.FindSpesificReligionTaskRequest(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific religion request task",
		})
	}

	userData, _ := handler.userUsecase.ReadSpecificUser(data.UserId)

	response := entity.UserReligionReqTaskCore{
		Id:          data.Id,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Title:       data.Title,
		Point:       data.Point,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Type:        data.Type,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get spesific request task",
		"data":    response,
	})
}

func (handler *TaskController) UploadReligionTaskRequest(e echo.Context) error {
	input := dto.UserReqReligionTaskRequest{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	userId, _, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	image, err := e.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return e.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "No file uploaded",
			})
		}
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error uploading file",
		})
	}

	dataInput := entity.UserReligionReqTaskCore{
		Title:       input.Title,
		Point:       input.Point,
		Image:       input.Image,
		Description: input.Description,
	}
	dataInput.UserId = userId

	errUpload := handler.taskUsecase.UploadReligionTaskRequest(dataInput, image)
	if errUpload != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error upload religion task request",
			"error":   errUpload.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes upload religion task request",
	})
}

func (handler *TaskController) GetAllUserReligionTaskRequest(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	data, err := handler.taskUsecase.GetAllUserReligionTaskRequest()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user religion request task",
		})
	}

	dataList := []entity.UserReligionReqTaskCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)

		result := entity.UserReligionReqTaskCore{
			Id:          v.Id,
			Title:       v.Title,
			Point:       v.Point,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Type:        v.Type,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user religion request task",
		"data":    dataList,
	})
}

func (handler *TaskController) UpdateTaskReligionReqStatus(e echo.Context) error {
	_, role, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	if role != "admin" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "access denied",
		})
	}

	idParams := e.Param("id")

	data := dto.UserReqReligionTaskRequest{}
	if errBind := e.Bind(&data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	//get userId & taskId
	dataTask, _ := handler.taskUsecase.FindSpesificReligionTaskRequest(idParams)

	status := entity.UserReligionReqTaskCore{
		UserId:      dataTask.UserId,
		UserName:    dataTask.UserName,
		Title:       data.Title,
		Image:       data.Image,
		Description: data.Description,
		Point:       data.Point,
		Status:      data.Status,
		Message:     data.Message,
	}

	errUpdate := handler.taskUsecase.UpdateTaskReligionReqStatus(idParams, status)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating religion request task status",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "religion task request status updated",
	})
}
