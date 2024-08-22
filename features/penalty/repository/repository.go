package repository

import (
	"errors"
	"tugaskita/features/penalty/entity"
	"tugaskita/features/penalty/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PenaltyRepository struct {
	db *gorm.DB
}

func NewPenaltyRepository(db *gorm.DB) entity.PenaltyDataInterface {
	return &PenaltyRepository{
		db: db,
	}
}

// CreatePenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) CreatePenalty(input entity.PenaltyCore) error {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	data := entity.PenaltyCoreToPenaltyModel(input)
	data.Id = newUUID
	tx := penaltyRepo.db.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeletePenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) DeletePenalty(id string) error {
	dataPenalty := model.Penalty{}

	tx := penaltyRepo.db.Where("id = ? ", id).Delete(&dataPenalty)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("penalty not found")
	}

	return nil
}

// FindAllPenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) FindAllPenalty() ([]entity.PenaltyCore, error) {
	var dataPenalty []model.Penalty

	errData := penaltyRepo.db.Find(&dataPenalty).Error
	if errData != nil {
		return nil, errData
	}

	dataResponse := make([]entity.PenaltyCore, len(dataPenalty))
	for i, v := range dataPenalty {
		dataResponse[i] = entity.PenaltyCore{
			Id:          v.Id,
			UserId:      v.UserId,
			Point:       v.Point,
			Description: v.Description,
			Date:        v.Date,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return dataResponse, nil
}

// FindSpecificPenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) FindSpecificPenalty(id string) (entity.PenaltyCore, error) {
	dataPenalty := model.Penalty{}

	tx := penaltyRepo.db.Where("id = ? ", id).First(&dataPenalty)
	if tx.Error != nil {
		return entity.PenaltyCore{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return entity.PenaltyCore{}, errors.New("penalty not found")
	}

	dataResponse := entity.PenaltyModelToPenaltyCore(dataPenalty)
	return dataResponse, nil
}

// UpdatePenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) UpdatePenalty(id string, data entity.PenaltyCore) error {
	dataPenalty := entity.PenaltyCoreToPenaltyModel(data)

	tx := penaltyRepo.db.Where("id = ?", id).Updates(&dataPenalty)
	if tx.Error != nil {
		if tx.Error != nil {
			return tx.Error
		}
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("penalty not found")
	}

	return nil
}

// FindAllPenaltyHistory implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) FindAllPenaltyHistory(id string) ([]entity.PenaltyCore, error) {
	var penalty []model.Penalty
	penaltyRepo.db.Where("user_id=?", id).Find(&penalty)

	datapPenalty := entity.ListPenaltyModelToListPenaltyCore(penalty)
	return datapPenalty, nil
}

// GetTotalPenalty implements entity.PenaltyDataInterface.
func (penaltyRepo *PenaltyRepository) GetTotalPenalty(id string) (int, error) {
	var count int64

	errPenalty := penaltyRepo.db.Model(&model.Penalty{}).
	Where("user_id = ?", id).
	Count(&count).Error
	if errPenalty != nil {
		return 0, errPenalty
	}

	return int(count), nil
}
