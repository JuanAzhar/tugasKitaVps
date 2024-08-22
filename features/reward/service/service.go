package service

import (
	"errors"
	"mime/multipart"
	"strconv"
	"tugaskita/features/reward/entity"
	user "tugaskita/features/user/entity"
)

type RewardService struct {
	RewardRepo entity.RewardDataInterface
	UserRepo   user.UserDataInterface
}

func NewRewardService(rewardRepo entity.RewardDataInterface, userRepo user.UserDataInterface) entity.RewardUseCaseInterface {
	return &RewardService{
		RewardRepo: rewardRepo,
		UserRepo:   userRepo,
	}
}

// CreateReward implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) CreateReward(input entity.RewardCore, image *multipart.FileHeader) error {

	if input.Name == "" {
		return errors.New("name and image can't be empty")
	}

	if input.Price < 0 || input.Stock < 0 {
		return errors.New("price and stock can't less then 0")
	}

	if image != nil && image.Size > 10*1024*1024 {
		return errors.New("image file size should be less than 10 MB")
	}

	err := rewardUC.RewardRepo.CreateReward(input, image)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) DeleteReward(rewardId string) error {
	if rewardId == "" {
		return errors.New("insert reward id")
	}

	_, err := rewardUC.RewardRepo.FindById(rewardId)
	if err != nil {
		return errors.New("reward not found")
	}

	errDelete := rewardUC.RewardRepo.DeleteReward(rewardId)
	if errDelete != nil {
		return errors.New("can't delete reward")
	}

	return nil
}

// FindAllReward implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) FindAllReward() ([]entity.RewardCore, error) {
	data, err := rewardUC.RewardRepo.FindAllReward()
	if err != nil {
		return nil, errors.New("error get data")
	}

	return data, nil
}

// FindById implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) FindById(rewardId string) (entity.RewardCore, error) {
	if rewardId == "" {
		return entity.RewardCore{}, errors.New("reward ID is required")
	}

	task, err := rewardUC.RewardRepo.FindById(rewardId)
	if err != nil {
		return entity.RewardCore{}, err
	}

	return task, nil
}

// UpdateReward implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) UpdateReward(rewardId string, data entity.RewardCore, image *multipart.FileHeader) error {
	if data.Name == "" {
		return errors.New("name can't be empty")
	}

	if data.Price < 0 || data.Stock < 0 {
		return errors.New("price and stock can't less then 0")
	}

	// Validasi ukuran file gambar jika gambar diunggah
	if image != nil {
		if image.Size > 10*1024*1024 {
			return errors.New("image file size should be less than 10 MB")
		}
	}

	err := rewardUC.RewardRepo.UpdateReward(rewardId, data, image)
	if err != nil {
		return err
	}

	return nil
}

// FindAllRewardHistory implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) FindAllRewardHistory(userId string) ([]entity.UserRewardRequestCore, error) {
	data, err := rewardUC.RewardRepo.FindAllRewardHistory(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindAllUploadReward implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) FindAllUploadReward() ([]entity.UserRewardRequestCore, error) {
	userReward, err := rewardUC.RewardRepo.FindAllUploadReward()
	if err != nil {
		return nil, errors.New("error get user reward request")
	}
	return userReward, nil
}

// UploadRewardRequest implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) UploadRewardRequest(input entity.UserRewardRequestCore) error {
	userData, errUser := rewardUC.UserRepo.ReadSpecificUser(input.UserId)
	if errUser != nil {
		return errors.New("failed get user")
	}

	userPoint, _ := strconv.Atoi(userData.TotalPoint)

	rewardData, errReward := rewardUC.RewardRepo.FindById(input.RewardId)
	if errReward != nil {
		return errors.New("failed get reward")
	}

	if rewardData.Stock < 1 {
		return errors.New("not enough stock")
	}

	totalPrice := rewardData.Price * input.Amount

	count := userPoint - totalPrice

	userData.TotalPoint = strconv.Itoa(count)
	input.TotalPrice = totalPrice
	input.Price = rewardData.Price

	if userPoint < totalPrice {
		return errors.New("not enough point")
	}

	//update user
	errUserUpdate := rewardUC.UserRepo.UpdatePoint(input.UserId, userData)
	if errUserUpdate != nil {
		return errors.New("failed update user point")
	}

	err := rewardUC.RewardRepo.UploadRewardRequest(input)
	if err != nil {
		return errors.New("failed request reward")
	}

	return nil
}

// FindUserRewardById implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) FindUserRewardById(id string) (entity.UserRewardRequestCore, error) {
	reward, err := rewardUC.RewardRepo.FindUserRewardById(id)
	if err != nil {
		return entity.UserRewardRequestCore{}, err
	}

	return reward, nil
}

// UpdateReqRewardStatus implements entity.RewardUseCaseInterface.
func (rewardUC *RewardService) UpdateReqRewardStatus(rewardId string, data entity.UserRewardRequestCore) error {
	// user data
	userData, errUser := rewardUC.UserRepo.ReadSpecificUser(data.UserId)
	if errUser != nil {
		return errors.New("failed get user")
	}

	//reward data
	rewardData, errReward := rewardUC.RewardRepo.FindById(data.RewardId)
	if errReward != nil {
		return errors.New("failed get reward")
	}

	//reward request
	rewardReqData, errRewardReq := rewardUC.RewardRepo.FindUserRewardById(rewardId)
	if errRewardReq != nil {
		return errors.New("failed get user reward request")
	}

	if rewardReqData.Status == "Diterima" {
		return errors.New("you already accept this request")
	}

	if rewardReqData.Status == "Ditolak" {
		return errors.New("you already reject this request")
	}

	if rewardData.Stock < 1 {
		return errors.New("not enough stock")
	}

	if data.Status == "Ditolak" {
		userPoint, _ := strconv.Atoi(userData.TotalPoint)
		count := userPoint + rewardReqData.TotalPrice
		userData.TotalPoint = strconv.Itoa(count)
	}

	//update user
	errUserUpdate := rewardUC.UserRepo.UpdatePoint(data.UserId, userData)
	if errUserUpdate != nil {
		return errors.New("failed update user point")
	}

	//save history
	if data.Status == "Diterima" {
		//update history
		amount := strconv.Itoa(rewardReqData.Amount)
		historyData := user.UserPointCore{
			UserId:   data.UserId,
			Type:     "Reward",
			Point:    rewardReqData.TotalPrice,
			TaskName: "Change " + amount + " " + rewardReqData.RewardName,
		}
		errUserHistory := rewardUC.UserRepo.PostUserPointHistory(historyData)
		if errUserHistory != nil {
			return errors.New("failed add user history point")
		}
	}

	//update status
	err := rewardUC.RewardRepo.UpdateReqRewardStatus(rewardId, data)
	if err != nil {
		return err
	}

	return nil
}
