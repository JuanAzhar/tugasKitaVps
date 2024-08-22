package entity

import "tugaskita/features/user/model"

func UserCoreToUserModel(data UserCore) model.Users {
	return model.Users{
		ID:         data.ID,
		Name:       data.Name,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Email:      data.Email,
		Password:   data.Password,
		Role:       data.Role,
		Point:      data.Point,
		TotalPoint: data.TotalPoint,
	}
}

func UserModelToUserCore(data model.Users) UserCore {
	return UserCore{
		ID:         data.ID,
		Name:       data.Name,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Email:      data.Email,
		Password:   data.Password,
		Role:       data.Role,
		Point:      data.Point,
		TotalPoint: data.TotalPoint,
	}
}

func UserPointCoreToUserPointModel(data UserPointCore) model.UserPoint {
	return model.UserPoint{
		Id:        data.Id,
		UserId:    data.UserId,
		Type:      data.Type,
		TaskName:  data.TaskName,
		Point:     data.Point,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func UserPointModelToUserPointCore(data model.UserPoint) UserPointCore {
	return UserPointCore{
		Id:        data.Id,
		UserId:    data.UserId,
		Type:      data.Type,
		TaskName:  data.TaskName,
		Point:     data.Point,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func ListUserPointModelToListUserPointCore(data []model.UserPoint) []UserPointCore {
	dataUser := []UserPointCore{}
	for _, v := range data {
		result := UserPointModelToUserPointCore(v)
		dataUser = append(dataUser, result)
	}
	return dataUser
}
