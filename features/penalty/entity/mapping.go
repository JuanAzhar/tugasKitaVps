package entity

import "tugaskita/features/penalty/model"

func PenaltyCoreToPenaltyModel(data PenaltyCore) model.Penalty {
	return model.Penalty{
		Id:          data.Id,
		UserId:      data.UserId,
		Point:       data.Point,
		Description: data.Description,
		Date:        data.Date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func PenaltyModelToPenaltyCore(data model.Penalty) PenaltyCore{
	return PenaltyCore{
		Id:          data.Id,
		UserId:      data.UserId,
		Point:       data.Point,
		Description: data.Description,
		Date:        data.Date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ListPenaltyModelToListPenaltyCore(data []model.Penalty) []PenaltyCore{
	dataPenalty := []PenaltyCore{}
	for _, v := range data {
		result := PenaltyModelToPenaltyCore(v)
		dataPenalty = append(dataPenalty, result)
	}
	return dataPenalty
}