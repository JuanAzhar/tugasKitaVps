package handler

import (
	"mime/multipart"
	"net/http"
	dto "tugaskita/features/user/dto"
	"tugaskita/features/user/entity"
	middleware "tugaskita/utils/jwt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUsecase entity.UserUseCaseInterface
}

func New(userUC entity.UserUseCaseInterface) *UserController {
	return &UserController{
		userUsecase: userUC,
	}
}

func (handler *UserController) Register(e echo.Context) error {
    input := dto.UserRequest{}
    errBind := e.Bind(&input)
    if errBind != nil {
        return e.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "error bind data: " + errBind.Error(),
        })
    }

    // Retrieve image file from the request, it can be nil if not provided
    image, err := e.FormFile("image")
    if err != nil && err != http.ErrMissingFile {
        return e.JSON(http.StatusBadRequest, map[string]interface{}{
            "message": "error uploading file",
        })
    }

    data := entity.UserCore{
        Name:     input.Name,
        Address:  input.Address,
        School:   input.School,
        Class:    input.Class,
        Religion: input.Religion,
        Email:    input.Email,
        Password: input.Password,
    }

    row, errUser := handler.userUsecase.Register(data, image)
    if errUser != nil {
        return e.JSON(http.StatusBadRequest, map[string]any{
            "message": "error creating user",
            "error":   errUser.Error(),
        })
    }

    return e.JSON(http.StatusOK, map[string]any{
        "message": "success creating user",
        "data":    row,
    })
}

func (handler *UserController) Login(e echo.Context) error {
	input := new(dto.UserRequest)
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	data := entity.UserCore{
		Email:    input.Email,
		Password: input.Password,
	}

	data, token, err := handler.userUsecase.Login(data.Email, data.Password)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error login",
			"error":   err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "login success",
		"email":   data.Email,
		"token":   token,
	})
}

func (handler *UserController) DeleteUser(e echo.Context) error {
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
	err := handler.userUsecase.DeleteUser(idParams)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error deleting user",
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

func (handler *UserController) ReadSpecificUser(e echo.Context) error {
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

	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "user not found",
		})
	}

	data, err := handler.userUsecase.ReadSpecificUser(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific user",
		})
	}

	response := dto.UserResponse{
		Id:         data.ID,
		Name:       data.Name,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Role:       data.Role,
		Religion:   data.Religion,
		Email:      data.Email,
		Image:      data.Image,
		Point:      data.Point,
		TotalPoint: data.TotalPoint,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get user",
		"data":    response,
	})
}

func (handler *UserController) ReadProfileUser(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	println("user Id : ", userId)

	idCheck, err := uuid.Parse(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "user not found",
		})
	}

	data, err := handler.userUsecase.ReadSpecificUser(idCheck.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get profile user",
		})
	}

	response := dto.UserResponse{
		Id:         data.ID,
		Name:       data.Name,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Image:      data.Image,
		Email:      data.Email,
		Religion:   data.Religion,
		Point:      data.Point,
		TotalPoint: data.TotalPoint,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get user profile",
		"data":    response,
	})
}

func (handler *UserController) ReadAllUser(e echo.Context) error {
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

	data, err := handler.userUsecase.ReadAllUser()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user",
		})
	}

	dataList := []dto.UserResponse{}
	for _, v := range data {
		result := dto.UserResponse{
			Id:         v.ID,
			Name:       v.Name,
			Address:    v.Address,
			School:     v.School,
			Class:      v.Class,
			Image:      v.Image,
			Email:      v.Email,
			Role:       v.Role,
			Religion:   v.Religion,
			Point:      v.Point,
			TotalPoint: v.TotalPoint,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user",
		"data":    dataList,
	})
}

func (handler *UserController) UpdateSiswa(e echo.Context) error {
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

	data := new(dto.UserRequest)
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

	// Membuat objek userData tanpa menyertakan gambar
	userData := entity.UserCore{
		Name:     data.Name,
		Email:    data.Email,
		Address:  data.Address,
		School:   data.School,
		Class:    data.Class,
		Password: data.Password,
		Religion: data.Religion,
		Point:    data.Point,
	}

	// Panggil fungsi UpdateSiswa di usecase, kirimkan image jika ada
	errUpdate := handler.userUsecase.UpdateSiswa(idParams, userData, image)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating user",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "user updated successfully",
	})
}

func (handler *UserController) GetRankUser(e echo.Context) error {
	data, err := handler.userUsecase.GetRankUser()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user",
		})
	}

	dataList := []dto.UserRankResponse{}
	for _, v := range data {
		result := dto.UserRankResponse{
			Name:  v.Name,
			Point: v.Point,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user rank",
		"data":    dataList,
	})
}

func (handler *UserController) ChangePassword(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data := new(dto.UserRequest)
	if errBind := e.Bind(data); errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error binding data",
		})
	}

	userData := entity.UserCore{
		Password: data.Password,
	}

	errUpdate := handler.userUsecase.ChangePassword(userId, userData)
	if errUpdate != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error updating password",
			"error":   errUpdate.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "password updated",
	})
}

func (handler *UserController) AnnualResetPoint(e echo.Context) error {
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

	errReset := handler.userUsecase.AnnualResetPoint()
	if errReset != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "failed reset point",
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "total point reset successfull",
	})
}

func (handler *UserController) MonthlyResetPoint(e echo.Context) error {
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

	errReset := handler.userUsecase.MonthlyResetPoint()
	if errReset != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "failed reset point",
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"message": "point reset successfull",
	})
}

func (handler *UserController) GetAllUserPointHistory(e echo.Context) error {
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

	data, err := handler.userUsecase.GetAllUserPointHistory()
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get all user point history",
		})
	}

	dataList := []entity.UserPointCore{}
	for _, v := range data {
		result := entity.UserPointCore{
			Id:        v.Id,
			UserId:    v.UserId,
			Type:      v.Type,
			TaskName:  v.TaskName,
			Point:     v.Point,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all user point history",
		"data":    dataList,
	})
}

func (handler *UserController) GetSpecificUserPointHistory(e echo.Context) error {
	idParamstr := e.Param("id")

	idParams, err := uuid.Parse(idParamstr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "history not found",
		})
	}

	data, err := handler.userUsecase.GetSpecificUserPointHistory(idParams.String())
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get specific history",
		})
	}

	response := entity.UserPointCore{
		Id:        data.Id,
		UserId:    data.UserId,
		Type:      data.Type,
		TaskName:  data.TaskName,
		Point:     data.Point,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get history",
		"data":    response,
	})
}

func (handler *UserController) GetUserPointHistory(e echo.Context) error {
	userId, _, _, err := middleware.ExtractTokenUserId(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	data, err := handler.userUsecase.GetUserPointHistory(userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error get history point",
		})
	}

	dataList := []entity.UserPointCore{}
	for _, v := range data {
		result := entity.UserPointCore{
			Id:        v.Id,
			UserId:    v.UserId,
			Type:      v.Type,
			TaskName:  v.TaskName,
			Point:     v.Point,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		dataList = append(dataList, result)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "get all point history",
		"data":    dataList,
	})
}

func (handler *UserController) PostUserPointHistory(e echo.Context) error {
	input := entity.UserPointCore{}
	errBind := e.Bind(&input)
	if errBind != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error bind data",
		})
	}

	errTask := handler.userUsecase.PostUserPointHistory(input)
	if errTask != nil {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"message": "error create history",
			"error":   errTask.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]any{
		"message": "succes create history",
	})
}
