package service

import (
	"errors"
	"strconv"
	"time"
	"tugaskita/features/penalty/entity"
	user "tugaskita/features/user/entity"
)

type PenaltyService struct {
	PenaltyRepo entity.PenaltyDataInterface
	UserRepo    user.UserDataInterface
}

func NewPenaltyService(penaltyRepo entity.PenaltyDataInterface, userRepo user.UserDataInterface) entity.PenaltyUseCaseInterface {
	return &PenaltyService{
		PenaltyRepo: penaltyRepo,
		UserRepo:    userRepo,
	}
}

// CreatePenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) CreatePenalty(input entity.PenaltyCore) error {
	if input.Description == "" || input.UserId == "" {
		return errors.New("description or userId can't empty")
	}

	if input.Point < 0 {
		return errors.New("point can't less then 0")
	}

	layout := "2006-01-02"

	_, errParse := time.Parse(layout, input.Date)
	if errParse != nil {
		return errors.New("date must be in 'yyyy-mm-dd'")
	}

	_, errUser := penaltyUC.UserRepo.ReadSpecificUser(input.UserId)
	if errUser != nil {
		return errors.New("user not found")
	}

	userData, errUser := penaltyUC.UserRepo.ReadSpecificUser(input.UserId)
	if errUser != nil {
		return errors.New("failed get user")
	}

	//reduce total point
	userTotalPoint, _ := strconv.Atoi(userData.TotalPoint)

	countTotal := userTotalPoint - input.Point

	userData.TotalPoint = strconv.Itoa(countTotal)

	//reduce point
	userPoint, _ := strconv.Atoi(userData.Point)

	count := userPoint - input.Point

	userData.Point = strconv.Itoa(count)
	//update user
	errUserUpdate := penaltyUC.UserRepo.UpdatePoint(input.UserId, userData)
	if errUserUpdate != nil {
		return errors.New("failed update user point")
	}

	//update history
	historyData := user.UserPointCore{
		UserId:   input.UserId,
		Type:     "Penalty",
		Point:    input.Point,
		TaskName: input.Description,
	}
	errUserHistory := penaltyUC.UserRepo.PostUserPointHistory(historyData)
	if errUserHistory != nil {
		return errors.New("failed add user history point")
	}

	//create penalty
	err := penaltyUC.PenaltyRepo.CreatePenalty(input)
	if err != nil {
		return err
	}

	return nil
}

// DeletePenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) DeletePenalty(id string) error {
	if id == "" {
		return errors.New("insert penalty id")
	}

	_, err := penaltyUC.PenaltyRepo.FindSpecificPenalty(id)
	if err != nil {
		return errors.New("penalty not found")
	}

	errDelete := penaltyUC.PenaltyRepo.DeletePenalty(id)
	if errDelete != nil {
		return errors.New("can't delete penalty")
	}

	return nil
}

// FindAllPenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) FindAllPenalty() ([]entity.PenaltyCore, error) {
	data, err := penaltyUC.PenaltyRepo.FindAllPenalty()
	if err != nil {
		return nil, errors.New("error get data")
	}

	return data, nil
}

// FindSpecificPenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) FindSpecificPenalty(id string) (entity.PenaltyCore, error) {
	if id == "" {
		return entity.PenaltyCore{}, errors.New("penalty ID is required")
	}

	task, err := penaltyUC.PenaltyRepo.FindSpecificPenalty(id)
	if err != nil {
		return entity.PenaltyCore{}, err
	}

	return task, nil
}

// UpdatePenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) UpdatePenalty(id string, data entity.PenaltyCore) error {
	layout := "2006-01-02"
	_, errParse := time.Parse(layout, data.Date)
	if errParse != nil {
		return errors.New("date must be in 'yyyy-mm-dd'")
	}

	//get user
	userData, errUser := penaltyUC.UserRepo.ReadSpecificUser(data.UserId)
	if errUser != nil {
		return errors.New("failed get user")
	}

	//get penalty
	penaltyData, _ := penaltyUC.PenaltyRepo.FindSpecificPenalty(id)

	userPoint, _ := strconv.Atoi(userData.Point)

	count := (userPoint + penaltyData.Point) - data.Point

	userData.Point = strconv.Itoa(count)

	//update user
	errUserUpdate := penaltyUC.UserRepo.UpdatePoint(data.UserId, userData)
	if errUserUpdate != nil {
		return errors.New("failed update user point")
	}

	err := penaltyUC.PenaltyRepo.UpdatePenalty(id, data)
	if err != nil {
		return err
	}

	return nil
}

// FindAllPenaltyHistory implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) FindAllPenaltyHistory(id string) ([]entity.PenaltyCore, error) {
	data, err := penaltyUC.PenaltyRepo.FindAllPenaltyHistory(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetTotalPenalty implements entity.PenaltyUseCaseInterface.
func (penaltyUC *PenaltyService) GetTotalPenalty(id string) (int, error) {
	penalty, err := penaltyUC.PenaltyRepo.GetTotalPenalty(id)
	if err != nil {
		return 0, errors.New("error count user penalty")
	}

	return penalty, nil
}