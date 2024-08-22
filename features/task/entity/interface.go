package entity

import (
	"mime/multipart"
)

type TaskDataInterface interface {
	CreateTask(input TaskCore) error
	FindAllTask() ([]TaskCore, error)
	FindById(taskId string) (TaskCore, error)
	UpdateTask(taskId string, data TaskCore) error
	DeleteTask(taskId string) error

	UpdateTaskStatus(taskId string, data UserTaskUploadCore) error
	UpdateTaskReqStatus(id string, data UserTaskSubmissionCore) error
	FindUserTaskById(id string) (UserTaskUploadCore, error)
	FindUserTaskReqById(id string) (UserTaskSubmissionCore, error)
	FindAllUserTask() ([]UserTaskUploadCore, error)

	UploadTask(input UserTaskUploadCore, image *multipart.FileHeader) error
	UploadTaskRequest(input UserTaskSubmissionCore, image *multipart.FileHeader) error
	FindAllRequestTask() ([]UserTaskSubmissionCore, error)
	FindAllClaimedTask(userId string) ([]UserTaskUploadCore, error)
	FindAllRequestTaskHistory(userId string) ([]UserTaskSubmissionCore, error)
	FindTasksNotClaimedByUser(userId string) ([]TaskCore, error)

	CountUserClearTask(id string) (int, error)

	CreateTaskReligion(input ReligionTaskCore) error
	FindAllTaskReligion() ([]ReligionTaskCore, error)
	FindByIdReligionTask(taskId string) (ReligionTaskCore, error)
	UpdateTaskReligion(taskId string, data ReligionTaskCore) error
	DeleteTaskReligion(taskId string) error
	FindTaskByDateAndReligion(date string, religion string) ([]ReligionTaskCore, error)
	FindTaskByDateAndReligionNon(date string, religion string) ([]ReligionTaskCore, error)

	UploadTaskReligion(input UserReligionTaskUploadCore, image *multipart.FileHeader) error
	FindAllReligionTaskUser(religion string, userId string) ([]ReligionTaskCore, error)
	FindAllReligionTaskHistory(userId string) ([]UserReligionTaskUploadCore, error)

	FindAllUserReligionTaskUpload() ([]UserReligionTaskUploadCore, error)
	FindSpecificUserReligionTaskUpload(userId string) (UserReligionTaskUploadCore, error)
	UpdateReligionTaskStatus(id string, data UserReligionTaskUploadCore) error

	UploadReligionTaskRequest(input UserReligionReqTaskCore, image *multipart.FileHeader) error
	FindAllReligionTaskRequestHistory(userId string) ([]UserReligionReqTaskCore, error)
	FindSpesificReligionTaskRequest(id string) (UserReligionReqTaskCore, error)

	GetAllUserReligionTaskRequest()([]UserReligionReqTaskCore, error)
	UpdateTaskReligionReqStatus(id string, data UserReligionReqTaskCore) error
}

type TaskUseCaseInterface interface {
	CreateTask(input TaskCore) error
	FindAllTask() ([]TaskCore, error)
	FindById(taskId string) (TaskCore, error)
	UpdateTask(taskId string, data TaskCore) error
	DeleteTask(taskId string) error

	UpdateTaskStatus(taskId string, data UserTaskUploadCore) error
	UpdateTaskReqStatus(id string, data UserTaskSubmissionCore) error
	FindUserTaskById(id string) (UserTaskUploadCore, error)
	FindUserTaskReqById(id string) (UserTaskSubmissionCore, error)
	FindAllUserTask() ([]UserTaskUploadCore, error)

	UploadTask(input UserTaskUploadCore, image *multipart.FileHeader) error
	UploadTaskRequest(input UserTaskSubmissionCore, image *multipart.FileHeader) error
	FindAllRequestTask() ([]UserTaskSubmissionCore, error)
	FindAllClaimedTask(userId string) ([]UserTaskUploadCore, error)
	FindAllRequestTaskHistory(userId string) ([]UserTaskSubmissionCore, error)
	FindTasksNotClaimedByUser(userId string) ([]TaskCore, error)

	CountUserClearTask(id string) (int, error)

	CreateTaskReligion(input ReligionTaskCore) error
	FindAllTaskReligion() ([]ReligionTaskCore, error)
	FindByIdReligionTask(taskId string) (ReligionTaskCore, error)
	UpdateTaskReligion(taskId string, data ReligionTaskCore) error
	DeleteTaskReligion(taskId string) error

	UploadTaskReligion(input UserReligionTaskUploadCore, image *multipart.FileHeader) error
	FindAllReligionTaskUser(religion string, userId string) ([]ReligionTaskCore, error)
	FindAllReligionTaskHistory(userId string) ([]UserReligionTaskUploadCore, error)

	FindAllUserReligionTaskUpload() ([]UserReligionTaskUploadCore, error)
	FindSpecificUserReligionTaskUpload(id string) (UserReligionTaskUploadCore, error)
	UpdateReligionTaskStatus(id string, data UserReligionTaskUploadCore) error

	UploadReligionTaskRequest(input UserReligionReqTaskCore, image *multipart.FileHeader) error
	FindAllReligionTaskRequestHistory(userId string) ([]UserReligionReqTaskCore, error)
	FindSpesificReligionTaskRequest(id string) (UserReligionReqTaskCore, error)

	GetAllUserReligionTaskRequest()([]UserReligionReqTaskCore, error)
	UpdateTaskReligionReqStatus(id string, data UserReligionReqTaskCore) error
}
