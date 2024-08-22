package handler

import (
	"mime/multipart"
	"net/http"
	"tugaskita/features/reward/dto"
	"tugaskita/features/reward/entity"
	user "tugaskita/features/user/entity"
	middleware "tugaskita/utils/jwt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RewardController struct {
	rewardUsecase entity.RewardUseCaseInterface
	userUsecase   user.UserUseCaseInterface
}

func New(rewardUC entity.RewardUseCaseInterface, userUC user.UserUseCaseInterface) *RewardController {
	return &RewardController{
		rewardUsecase: rewardUC,
		userUsecase:   userUC,
	}
}

func (handler *RewardController) AddReward(e echo.Context) error {
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

	input := dto.RewardRequest{}
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

	data := entity.RewardCore{
		Name:  input.Name,
		Stock: input.Stock,
		Price: input.Price,
		Image: input.Image,
	}

	errTask := handler.rewardUsecase.CreateReward(data, image)
	if errTask != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error create reward",
			"error":   errTask.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes create reward",
	})
}

func (handler *RewardController) ReadAllReward(e echo.Context) error {
	data, err := handler.rewardUsecase.FindAllReward()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all reward",
		})
	}

	dataList := []entity.RewardCore{}
	for _, v := range data {
		result := entity.RewardCore{
			ID:        v.ID,
			Name:      v.Name,
			Stock:     v.Stock,
			Price:     v.Price,
			Image:     v.Image,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all reward",
		"data":    dataList,
	})
}

func (handler *RewardController) ReadSpecificReward(e echo.Context) error {

	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "reward not found",
		})
	}

	data, err := handler.rewardUsecase.FindById(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific reward",
		})
	}

	response := entity.RewardCore{
		ID:        data.ID,
		Name:      data.Name,
		Stock:     data.Stock,
		Price:     data.Price,
		Image:     data.Image,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get reward",
		"data":    response,
	})
}

func (handler *RewardController) DeleteReward(e echo.Context) error {
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
	err := handler.rewardUsecase.DeleteReward(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error deleting reward",
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "Reward deleted successfully",
	})
}

func (handler *RewardController) UpdateReward(e echo.Context) error {
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

	data := new(dto.RewardRequest)
	if errBind := e.Bind(data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	// Menginisialisasi variabel untuk file gambar
	var image *multipart.FileHeader
	image, err = e.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error uploading file",
		})
	}

	rewardData := entity.RewardCore{
		Name:  data.Name,
		Stock: data.Stock,
		Price: data.Price,
	}

	errUpdate := handler.rewardUsecase.UpdateReward(idParams, rewardData, image)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating reward",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "reward updated successfully",
	})
}

func (handler *RewardController) FindAllRewardHistory(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.rewardUsecase.FindAllRewardHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get history reward",
		})
	}

	dataList := []entity.UserRewardRequestCore{}
	for _, v := range data {

		rewardData, _ := handler.rewardUsecase.FindById(v.RewardId)

		result := entity.UserRewardRequestCore{
			Id:         v.Id,
			RewardId:   v.RewardId,
			RewardName: rewardData.Name,
			Price:      rewardData.Price,
			UserId:     v.UserId,
			Status:     v.Status,
			TotalPrice: v.TotalPrice,
			Amount:     v.Amount,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all reward history",
		"data":    dataList,
	})
}

func (handler *RewardController) FindAllUploadReward(e echo.Context) error {
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

	data, err := handler.rewardUsecase.FindAllUploadReward()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user reward",
		})
	}

	dataList := []entity.UserRewardRequestCore{}
	for _, v := range data {

		userData, _ := handler.userUsecase.ReadSpecificUser(v.UserId)
		rewardData, _ := handler.rewardUsecase.FindById(v.RewardId)

		result := entity.UserRewardRequestCore{
			Id:         v.Id,
			RewardId:   v.RewardId,
			RewardName: rewardData.Name,
			Price:      rewardData.Price,
			UserId:     v.UserId,
			Amount:     v.Amount,
			TotalPrice: v.TotalPrice,
			UserName:   userData.Name,
			Status:     v.Status,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user reward",
		"data":    dataList,
	})
}

func (handler *RewardController) UploadRewardRequest(e echo.Context) error {
	input := new(dto.RewardReqRequest)
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

	dataInput := entity.UserRewardRequestCore{
		RewardId: input.RewardId,
		UserId:   userId,
		Amount:   input.Amount,
	}

	err := handler.rewardUsecase.UploadRewardRequest(dataInput)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error upload request reward",
			"error":   err.Error(),
		})
	}
	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes upload request reward",
	})
}

func (handler *RewardController) FindUserRewardById(e echo.Context) error {
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

	data, err := handler.rewardUsecase.FindUserRewardById(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific reward",
		})
	}

	userData, _ := handler.userUsecase.ReadSpecificUser(data.UserId)
	rewardData, _ := handler.rewardUsecase.FindById(data.RewardId)

	response := entity.UserRewardRequestCore{
		Id:         data.Id,
		RewardId:   data.RewardId,
		RewardName: rewardData.Name,
		Price:      data.Price,
		UserId:     data.UserId,
		UserName:   userData.Name,
		Amount:     data.Amount,
		TotalPrice: data.TotalPrice,
		Status:     data.Status,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get reward",
		"data":    response,
	})
}

func (handler *RewardController) UpdateReqRewardStatus(e echo.Context) error {
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

	data := dto.RewardReqUpdateRequest{}
	if errBind := e.Bind(&data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	//get userId & rewardId
	rewardData, _ := handler.rewardUsecase.FindUserRewardById(idParams)

	status := entity.UserRewardRequestCore{
		RewardId: rewardData.RewardId,
		UserId:   rewardData.UserId,
		Status:   data.Status,
	}

	errUpdate := handler.rewardUsecase.UpdateReqRewardStatus(idParams, status)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating reward status",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "reward updated successfully",
	})
}
