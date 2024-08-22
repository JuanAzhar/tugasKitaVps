package repository

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"tugaskita/features/reward/entity"
	"tugaskita/features/reward/model"
	user "tugaskita/features/user/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RewardRepository struct {
	db             *gorm.DB
	userRepository user.UserDataInterface
}

func NewRewardRepository(db *gorm.DB, userRepository user.UserDataInterface) entity.RewardDataInterface {
	return &RewardRepository{
		db:             db,
		userRepository: userRepository,
	}
}

// CreateReward implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) CreateReward(input entity.RewardCore, image *multipart.FileHeader) error {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	file, err := image.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Define the directory where you want to save the image
	saveDir := "public/images/reward"
	os.MkdirAll(saveDir, os.ModePerm)

	// Define the file path
	filePath := filepath.Join(saveDir, image.Filename)

	// Replace backslashes with forward slashes for consistency
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded file data to the destination file
	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	input.Image = filePath

	data := entity.RewardCoreToRewardModel(input)
	data.ID = newUUID
	tx := rewardRepo.db.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeleteTask implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) DeleteReward(rewardId string) error {
	dataReward := model.Reward{}

	tx := rewardRepo.db.Where("id = ? ", rewardId).Delete(&dataReward)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("reward not found")
	}

	return nil
}

// FindAllReward implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) FindAllReward() ([]entity.RewardCore, error) {
	var reward []model.Reward
	rewardRepo.db.Find(&reward)

	dataReward := entity.ListRewardModelToRewardCore(reward)
	return dataReward, nil
}

// FindById implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) FindById(rewardId string) (entity.RewardCore, error) {
	dataReward := model.Reward{}

	tx := rewardRepo.db.Where("id = ? ", rewardId).First(&dataReward)
	if tx.Error != nil {
		return entity.RewardCore{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return entity.RewardCore{}, errors.New("reward not found")
	}

	dataResponse := entity.RewardModelToRewardCore(dataReward)
	return dataResponse, nil
}

// UpdateReward implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) UpdateReward(rewardId string, data entity.RewardCore, image *multipart.FileHeader) error {
	dataReward := entity.RewardCoreToRewardModel(data)

	if image != nil {
		file, err := image.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// Define the directory where you want to save the image
		saveDir := "public/images/reward"
		os.MkdirAll(saveDir, os.ModePerm)

		// Define the file path
		filePath := filepath.Join(saveDir, image.Filename)

		// Replace backslashes with forward slashes for consistency
		filePath = strings.ReplaceAll(filePath, "\\", "/")

		// Create the destination file
		dst, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination file
		if _, err = io.Copy(dst, file); err != nil {
			return err
		}

		dataReward.Image = filePath
	}

	tx := rewardRepo.db.Where("id = ?", rewardId).Updates(&dataReward)
	if tx.Error != nil {
		if tx.Error != nil {
			return tx.Error
		}
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("reward not found")
	}

	return nil
}

// FindAllUploadReward implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) FindAllUploadReward() ([]entity.UserRewardRequestCore, error) {
	var reward []model.UserRewardRequest

	errData := rewardRepo.db.Find(&reward).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserRewardRequestCore, len(reward))
	for i, v := range reward {
		mapData[i] = entity.UserRewardRequestCore{
			Id:         v.Id,
			RewardId:   v.RewardId,
			Price:      v.Price,
			TotalPrice: v.TotalPrice,
			Amount:     v.Amount,
			UserId:     v.UserId,
			Status:     v.Status,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		}
	}
	return mapData, nil
}

// UploadRewardRequest implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) UploadRewardRequest(input entity.UserRewardRequestCore) error {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	var inputData = model.UserRewardRequest{
		Id:         newUUID,
		RewardId:   input.RewardId,
		UserId:     input.UserId,
		Status:     input.Status,
		Amount:     input.Amount,
		TotalPrice: input.TotalPrice,
		Price:      input.Price,
	}

	errUpload := rewardRepo.db.Save(&inputData)
	if errUpload != nil {
		return errUpload.Error
	}

	return nil
}

// FindAllRewardRequestUser implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) FindAllRewardHistory(userId string) ([]entity.UserRewardRequestCore, error) {
	var reward []model.UserRewardRequest
	rewardRepo.db.Where("user_id=?", userId).Find(&reward)

	dataReward := entity.ListRewardUserModelToListRewardUserCore(reward)
	return dataReward, nil
}

// FindUserRewardById implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) FindUserRewardById(id string) (entity.UserRewardRequestCore, error) {
	var data model.UserRewardRequest

	errData := rewardRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserRewardRequestCore{}, errData
	}

	rewardData, _ := rewardRepo.FindById(data.RewardId)
	userData, _ := rewardRepo.userRepository.ReadSpecificUser(data.UserId)

	userCore := entity.UserRewardRequestCore{
		Id:         data.Id,
		RewardId:   data.RewardId,
		RewardName: rewardData.Name,
		UserName:   userData.Name,
		UserId:     data.UserId,
		Status:     data.Status,
		Amount:     data.Amount,
		Price:      data.Price,
		TotalPrice: data.TotalPrice,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}

	return userCore, nil
}

// UpdateReqRewardStatus implements entity.RewardDataInterface.
func (rewardRepo *RewardRepository) UpdateReqRewardStatus(rewardId string, data entity.UserRewardRequestCore) error {
	rewardData := entity.RewardUserCoreToRewardUserModel(data)

	//update status
	tx := rewardRepo.db.Where("id=?", rewardId).Updates(rewardData)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
