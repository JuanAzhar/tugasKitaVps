package entity

import (
	"tugaskita/features/task/model"
)

func TaskCoreToTaskModel(data TaskCore) model.Task {
	return model.Task{
		ID:          data.ID,
		AdminId:     data.AdminId,
		Title:       data.Title,
		Description: data.Description,
		Point:       data.Point,
		Message:     data.Message,
		Status:      data.Status,
		Type:        data.Type,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

}

func TaskModelToTaskCore(data model.Task) TaskCore {
	return TaskCore{
		ID:          data.ID,
		AdminId:     data.AdminId,
		Title:       data.Title,
		Description: data.Description,
		Point:       data.Point,
		Message:     data.Message,
		Status:      data.Status,
		Type:        data.Type,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

}

func ListTaskModelToTaskCore(data []model.Task) []TaskCore {
	dataTask := []TaskCore{}
	for _, v := range data {
		result := TaskModelToTaskCore(v)
		dataTask = append(dataTask, result)
	}
	return dataTask
}

func TaskUserModelToTaskUserCore(data model.UserTaskUpload) UserTaskUploadCore {
	return UserTaskUploadCore{
		Id:          data.Id,
		TaskId:      data.TaskId,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

}

func TaskUserCoreToTaskUserModel(data UserTaskUploadCore) model.UserTaskUpload {
	return model.UserTaskUpload{
		Id:          data.Id,
		TaskId:      data.TaskId,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ListTaskUserModelToTaskUserCore(data []model.UserTaskUpload) []UserTaskUploadCore {
	dataTask := []UserTaskUploadCore{}
	for _, v := range data {
		result := TaskUserModelToTaskUserCore(v)
		dataTask = append(dataTask, result)
	}
	return dataTask
}

func TaskUserReqModelToTaskUserReqCore(data model.UserTaskSubmission) UserTaskSubmissionCore {
	return UserTaskSubmissionCore{
		Id:          data.Id,
		Title:       data.Title,
		Point:       data.Point,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func TaskUserReqCoreToTaskUserReqModel(data UserTaskSubmissionCore) model.UserTaskSubmission {
	return model.UserTaskSubmission{
		Id:          data.Id,
		Title:       data.Title,
		Point:       data.Point,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ListTaskUserReqModelToTaskUserReqCore(data []model.UserTaskSubmission) []UserTaskSubmissionCore {
	dataTask := []UserTaskSubmissionCore{}
	for _, v := range data {
		result := TaskUserReqModelToTaskUserReqCore(v)
		dataTask = append(dataTask, result)
	}
	return dataTask
}

func ReligionTaskCoreToTaskModel(data ReligionTaskCore) model.ReligionTask {
	return model.ReligionTask{
		Id:          data.Id,
		Title:       data.Title,
		Description: data.Description,
		Type:        data.Type,
		Religion:    data.Religion,
		Point:       data.Point,
		Start_date:  data.Start_date,
		End_date:    data.End_date,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.CreatedAt,
	}

}

func ReligionTaskModelToTaskCore(data model.ReligionTask) ReligionTaskCore {
	return ReligionTaskCore{
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

}

func ListReligionTaskModelToReligionTaskCore(data []model.ReligionTask) []ReligionTaskCore {
	dataTask := []ReligionTaskCore{}
	for _, v := range data {
		result := ReligionTaskModelToTaskCore(v)
		dataTask = append(dataTask, result)
	}
	return dataTask
}

func ReligionTaskUploadCoreToReligionTaskUploadModel(data UserReligionTaskUploadCore) model.UserReligionTaskUpload {
	return model.UserReligionTaskUpload{
		Id:          data.Id,
		TaskId:      data.TaskId,
		UserId:      data.UserId,
		Image:       data.Image,
		Type:        data.Type,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ReligionTaskUploadModelToReligionTaskUploadCore(data model.UserReligionTaskUpload) UserReligionTaskUploadCore {
	return UserReligionTaskUploadCore{
		Id:          data.Id,
		TaskId:      data.TaskId,
		UserId:      data.UserId,
		Image:       data.Image,
		Type:        data.Type,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ReligionTaskReqModelToReligionTaskReqCore(data model.UserReligionReqTask) UserReligionReqTaskCore {
	return UserReligionReqTaskCore{
		Id:          data.Id,
		Title:       data.Title,
		Point:       data.Point,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ReligionTaskReqCoreToReligioinTaskReqModel(data UserReligionReqTaskCore) model.UserReligionReqTask {
	return model.UserReligionReqTask{
		Id:          data.Id,
		Title:       data.Title,
		Point:       data.Point,
		UserId:      data.UserId,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func ListReligionReqTaskModelToReligionReqTaskCore(data []model.UserReligionReqTask) []UserReligionReqTaskCore {
	dataTask := []UserReligionReqTaskCore{}
	for _, v := range data {
		result := ReligionTaskReqModelToReligionTaskReqCore(v)
		dataTask = append(dataTask, result)
	}
	return dataTask
}
