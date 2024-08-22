package handler

import (
	"net/http"
	"tugaskita/features/penalty/dto"
	"tugaskita/features/penalty/entity"
	user "tugaskita/features/user/entity"
	middleware "tugaskita/utils/jwt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PenaltyController struct {
	penaltyUsecase entity.PenaltyUseCaseInterface
	userUsecase    user.UserUseCaseInterface
}

func New(penaltyUC entity.PenaltyUseCaseInterface, userUC user.UserUseCaseInterface) *PenaltyController {
	return &PenaltyController{
		penaltyUsecase: penaltyUC,
		userUsecase:    userUC,
	}
}

func (handler *PenaltyController) CreatePenalty(e echo.Context) error {
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

	input := dto.PenaltyRequest{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	data := entity.PenaltyCore{
		UserId:      input.UserId,
		Description: input.Description,
		Point:       input.Point,
		Date:        input.Date,
	}

	errTask := handler.penaltyUsecase.CreatePenalty(data)
	if errTask != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error create penalty",
			"error":   errTask.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes create penalty",
	})
}

func (handler *PenaltyController) DeletePenalty(e echo.Context) error {
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
	err := handler.penaltyUsecase.DeletePenalty(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error deleting penalty",
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "penalty deleted successfully",
	})
}

func (handler *PenaltyController) FindAllPenalty(e echo.Context) error {
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

	data, err := handler.penaltyUsecase.FindAllPenalty()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all penalty user",
		})
	}

	dataList := []entity.PenaltyCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)

		result := entity.PenaltyCore{
			Id:          v.Id,
			UserId:      v.UserId,
			Description: v.Description,
			UserName: userData.Name,
			Point:       v.Point,
			Date:        v.Date,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all penalty user",
		"data":    dataList,
	})
}

func (handler *PenaltyController) FindSpecificPenalty(e echo.Context) error {
	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "penalty data not found",
		})
	}

	data, err := handler.penaltyUsecase.FindSpecificPenalty(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific penalty",
		})
	}

	userData, errData := handler.userUsecase.ReadSpecificUser(data.UserId)
	if errData != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific user",
		})
	}

	response := entity.PenaltyCore{
		Id:          data.Id,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Description: data.Description,
		Point:       data.Point,
		Date:        data.Date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get specific penalty",
		"data":    response,
	})
}

func (handler *PenaltyController) UpdatePenalty(e echo.Context) error {
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

	data := new(dto.PenaltyRequest)
	if errBind := e.Bind(data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	rewardData := entity.PenaltyCore{
		UserId:      data.UserId,
		Description: data.Description,
		Point:       data.Point,
		Date:        data.Date,
	}

	errUpdate := handler.penaltyUsecase.UpdatePenalty(idParams, rewardData)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating penalty",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "penalty updated successfully",
	})
}

func (handler *PenaltyController) FindAllPenaltyHistory(e echo.Context) error {
	userId, _, _, errRole := middleware.ExtractTokenUserId(e)
	if errRole != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": errRole.Error(),
		})
	}

	data, err := handler.penaltyUsecase.FindAllPenaltyHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all penalty history",
		})
	}

	dataList := []entity.PenaltyCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)

		result := entity.PenaltyCore{
			Id:          v.Id,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Description: v.Description,
			Point:       v.Point,
			Date:        v.Date,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all penalty history",
		"data":    dataList,
	})
}

func (handler *PenaltyController) CountUserPenalty(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	count, err := handler.penaltyUsecase.GetTotalPenalty(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get penalty sum " + err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all task cleared sum",
		"count":   count,
	})
}