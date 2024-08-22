package repository

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tugaskita/features/task/entity"
	"tugaskita/features/task/model"
	user "tugaskita/features/user/entity"
	userModel "tugaskita/features/user/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db             *gorm.DB
	userRepository user.UserDataInterface
}

func NewTaskRepository(db *gorm.DB, userRepository user.UserDataInterface) entity.TaskDataInterface {
	return &TaskRepository{
		db:             db,
		userRepository: userRepository,
	}
}

// CreateTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) CreateTask(input entity.TaskCore) error {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	data := entity.TaskCoreToTaskModel(input)
	data.ID = newUUID
	tx := taskRepo.db.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeleteTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) DeleteTask(taskId string) error {
	dataTask := model.Task{}

	tx := taskRepo.db.Where("id = ? ", taskId).Delete(&dataTask)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// FindAllMission implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllTask() ([]entity.TaskCore, error) {
	var task []model.Task
	taskRepo.db.Find(&task)

	dataTask := entity.ListTaskModelToTaskCore(task)
	return dataTask, nil
}

// FindById implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindById(taskId string) (entity.TaskCore, error) {
	dataTask := model.Task{}

	tx := taskRepo.db.Where("id = ? ", taskId).First(&dataTask)
	if tx.Error != nil {
		return entity.TaskCore{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return entity.TaskCore{}, errors.New("task not found")
	}

	dataResponse := entity.TaskModelToTaskCore(dataTask)
	return dataResponse, nil
}

// UpdateTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateTask(taskId string, data entity.TaskCore) error {
	dataTask := entity.TaskCoreToTaskModel(data)

	tx := taskRepo.db.Where("id = ?", taskId).Updates(&dataTask)
	if tx.Error != nil {
		if tx.Error != nil {
			return tx.Error
		}
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// UpdateTaskStatus implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateTaskStatus(taskId string, data entity.UserTaskUploadCore) error {
	var pointTask model.Task
	var userData userModel.Users
	taskData := entity.TaskUserCoreToTaskUserModel(data)

	// get task data
	errData := taskRepo.db.Where("id=?", data.TaskId).First(&pointTask).Error
	if errData != nil {
		return errData
	}

	// get user
	errUser := taskRepo.db.Where("id=?", data.UserId).First(&userData).Error
	if errUser != nil {
		return errUser
	}

	// update status
	tx := taskRepo.db.Where("id=?", taskId).Updates(taskData)
	if tx.Error != nil {
		return tx.Error
	}

	if taskData.Status == "Diterima" {
		//update user
		userPoint, _ := strconv.Atoi(userData.Point)
		userTotalPoint, _ := strconv.Atoi(userData.TotalPoint)

		count := userPoint + pointTask.Point
		countTotal := userTotalPoint + pointTask.Point

		userData.Point = strconv.Itoa(count)
		userData.TotalPoint = strconv.Itoa(countTotal)

		saveUser := user.UserModelToUserCore(userData)

		updateUser := taskRepo.userRepository.UpdatePoint(data.UserId, saveUser)
		if updateUser != nil {
			return updateUser
		}

		//update history
		historyData := user.UserPointCore{
			UserId:   data.UserId,
			Type:     "Task",
			Point:    pointTask.Point,
			TaskName: pointTask.Title,
		}
		errUserHistory := taskRepo.userRepository.PostUserPointHistory(historyData)
		if errUserHistory != nil {
			return errors.New("failed add user history point")
		}
	}

	return nil
}

// UpdateTaskReqStatus implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateTaskReqStatus(id string, data entity.UserTaskSubmissionCore) error {
	var pointTask model.UserTaskSubmission
	var userData userModel.Users
	taskData := entity.TaskUserReqCoreToTaskUserReqModel(data)

	// get user
	errUser := taskRepo.db.Where("id=?", data.UserId).First(&userData).Error
	if errUser != nil {
		return errUser
	}

	// update status
	tx := taskRepo.db.Where("id=?", id).Updates(taskData)
	if tx.Error != nil {
		return tx.Error
	}

	// get task data
	errData := taskRepo.db.Where("id=?", id).First(&pointTask).Error
	if errData != nil {
		return errData
	}

	if taskData.Status == "Diterima" {
		userPoint, _ := strconv.Atoi(userData.Point)
		userTotalPoint, _ := strconv.Atoi(userData.TotalPoint)

		count := userPoint + pointTask.Point
		countTotal := userTotalPoint + pointTask.Point

		userData.Point = strconv.Itoa(count)
		userData.TotalPoint = strconv.Itoa(countTotal)

		saveUser := user.UserModelToUserCore(userData)

		updateUser := taskRepo.userRepository.UpdatePoint(data.UserId, saveUser)
		if updateUser != nil {
			return updateUser
		}

		//update history
		historyData := user.UserPointCore{
			UserId:   data.UserId,
			Type:     "Submission",
			Point:    pointTask.Point,
			TaskName: pointTask.Title,
		}
		errUserHistory := taskRepo.userRepository.PostUserPointHistory(historyData)
		if errUserHistory != nil {
			return errors.New("failed add user history point")
		}
	}

	return nil
}

// FindAllClaimedTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllClaimedTask(userId string) ([]entity.UserTaskUploadCore, error) {
	var task []model.UserTaskUpload
	taskRepo.db.Where("user_id=?", userId).Find(&task)

	dataTask := make([]entity.UserTaskUploadCore, len(task))
	for i, v := range task {
		dataTask[i] = entity.UserTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			UserId:      v.UserId,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return dataTask, nil
}

// FindAllRequestTaskHistory implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllRequestTaskHistory(userId string) ([]entity.UserTaskSubmissionCore, error) {
	var task []model.UserTaskSubmission
	taskRepo.db.Where("user_id=?", userId).Find(&task)

	mapData := make([]entity.UserTaskSubmissionCore, len(task))
	for i, v := range task {

		userData, _ := taskRepo.userRepository.ReadSpecificUser(userId)
		mapData[i] = entity.UserTaskSubmissionCore{
			Id:          v.Id,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Title:       v.Title,
			Type:        v.Type,
			Point:       v.Point,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// FindAllTaskNotClaimed implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindTasksNotClaimedByUser(userId string) ([]entity.TaskCore, error) {
	var tasks []model.Task

	currentDate := time.Now().Format("2006-01-02")

	taskRepo.db.Raw(`
	SELECT * FROM tasks 
	WHERE id NOT IN (
		SELECT task_id FROM user_task_uploads 
		WHERE user_id = ? AND status != 'Ditolak'
	) AND status = 'Active' AND end_date >= ?
`, userId, currentDate).Scan(&tasks)

	data := entity.ListTaskModelToTaskCore(tasks)
	return data, nil
}

// UploadTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UploadTask(input entity.UserTaskUploadCore, image *multipart.FileHeader) error {
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
	saveDir := "public/images/uploadTask"
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

	var inputData = model.UserTaskUpload{
		Id:          newUUID,
		TaskId:      input.TaskId,
		UserId:      input.UserId,
		Image:       input.Image,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	errUpload := taskRepo.db.Save(&inputData)
	if errUpload != nil {
		return errUpload.Error
	}

	return nil
}

// FindUserTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllUserTask() ([]entity.UserTaskUploadCore, error) {
	var userTask []model.UserTaskUpload

	errData := taskRepo.db.Find(&userTask).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserTaskUploadCore, len(userTask))
	for i, v := range userTask {
		mapData[i] = entity.UserTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			UserId:      v.UserId,
			Image:       v.Image,
			Type:        v.Type,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// FindUserTaskById implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindUserTaskById(id string) (entity.UserTaskUploadCore, error) {
	var data model.UserTaskUpload

	errData := taskRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserTaskUploadCore{}, errData
	}

	userData, _ := taskRepo.userRepository.ReadSpecificUser(data.UserId)
	taskData, _ := taskRepo.FindById(data.TaskId)

	userCore := entity.UserTaskUploadCore{
		Id:          data.Id,
		TaskId:      data.TaskId,
		TaskName:    taskData.Title,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Type:        data.Type,
		Message:     data.Message,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return userCore, nil
}

// FindUserTaskReqById implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindUserTaskReqById(id string) (entity.UserTaskSubmissionCore, error) {
	var data model.UserTaskSubmission

	errData := taskRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserTaskSubmissionCore{}, errData
	}

	userData, _ := taskRepo.userRepository.ReadSpecificUser(data.UserId)

	userCore := entity.UserTaskSubmissionCore{
		Id:          data.Id,
		Title:       data.Title,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Image:       data.Image,
		Point:       data.Point,
		Description: data.Description,
		Status:      data.Status,
		Message:     data.Message,
		Type:        data.Type,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return userCore, nil
}

// UploadTaskRequest implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UploadTaskRequest(input entity.UserTaskSubmissionCore, image *multipart.FileHeader) error {
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
	saveDir := "public/images/uploadTaskRequest"
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

	var inputData = model.UserTaskSubmission{
		Id:          newUUID,
		UserId:      input.UserId,
		Title:       input.Title,
		Point:       input.Point,
		Image:       input.Image,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	errUpload := taskRepo.db.Save(&inputData)
	if errUpload != nil {
		return errUpload.Error
	}

	return nil
}

// FindAllRequestTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllRequestTask() ([]entity.UserTaskSubmissionCore, error) {
	var userTask []model.UserTaskSubmission

	errData := taskRepo.db.Find(&userTask).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserTaskSubmissionCore, len(userTask))
	for i, v := range userTask {
		mapData[i] = entity.UserTaskSubmissionCore{
			Id:          v.Id,
			UserId:      v.UserId,
			Title:       v.Title,
			Image:       v.Image,
			Description: v.Description,
			Point:       v.Point,
			Status:      v.Status,
			Message:     v.Message,
			Type:        v.Type,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// CountUserClearTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) CountUserClearTask(id string) (int, error) {
	var countUpload int64
	var countSubmission int64
	var countTask int64

	// jumlah pada tabel UserTaskUpload
	errUpload := taskRepo.db.Model(&model.UserTaskUpload{}).
		Where("user_id = ? AND status = ?", id, "Diterima").
		Count(&countUpload).Error
	if errUpload != nil {
		return 0, errUpload
	}

	// jumlah pada tabel UserTaskSubmission
	errSubmission := taskRepo.db.Model(&model.UserTaskSubmission{}).
		Where("user_id = ? AND status = ?", id, "Diterima").
		Count(&countSubmission).Error
	if errSubmission != nil {
		return 0, errSubmission
	}

	// jumlah pada tabel UserReligionTaskUpload
	errTask := taskRepo.db.Model(&model.UserReligionTaskUpload{}).
		Where("user_id = ? AND status = ?", id, "Diterima").
		Count(&countTask).Error
	if errTask != nil {
		return 0, errTask
	}

	// Menggabungkan hasil dari kedua tabel
	totalCount := int(countUpload + countSubmission + countTask)

	return totalCount, nil
}

// CreateTaskReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) CreateTaskReligion(input entity.ReligionTaskCore) error {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	data := entity.ReligionTaskCoreToTaskModel(input)
	data.Id = newUUID
	tx := taskRepo.db.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeleteTaskReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) DeleteTaskReligion(taskId string) error {
	dataTask := model.ReligionTask{}

	tx := taskRepo.db.Where("id = ? ", taskId).Delete(&dataTask)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// FindAllTaskReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllTaskReligion() ([]entity.ReligionTaskCore, error) {
	var task []model.ReligionTask
	taskRepo.db.Find(&task)

	dataTask := entity.ListReligionTaskModelToReligionTaskCore(task)
	return dataTask, nil
}

// FindByIdReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindByIdReligionTask(taskId string) (entity.ReligionTaskCore, error) {
	dataTask := model.ReligionTask{}

	tx := taskRepo.db.Where("id = ? ", taskId).First(&dataTask)
	if tx.Error != nil {
		return entity.ReligionTaskCore{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return entity.ReligionTaskCore{}, errors.New("task not found")
	}

	dataResponse := entity.ReligionTaskModelToTaskCore(dataTask)
	return dataResponse, nil
}

// UpdateTaskReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateTaskReligion(taskId string, data entity.ReligionTaskCore) error {
	dataTask := entity.ReligionTaskCoreToTaskModel(data)

	tx := taskRepo.db.Where("id = ?", taskId).Updates(&dataTask)
	if tx.Error != nil {
		if tx.Error != nil {
			return tx.Error
		}
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// FindTaskByDateAndReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindTaskByDateAndReligion(date string, religion string) ([]entity.ReligionTaskCore, error) {
	var tasks []model.ReligionTask
	err := taskRepo.db.Where("start_date = ? AND religion = ? AND title = ?", date, religion, "Subuh").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	var taskCores []entity.ReligionTaskCore
	for _, task := range tasks {
		taskCores = append(taskCores, entity.ReligionTaskModelToTaskCore(task))
	}
	return taskCores, nil
}

// FindTaskByDateAndReligionNon implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindTaskByDateAndReligionNon(date string, religion string) ([]entity.ReligionTaskCore, error) {
	var tasks []model.ReligionTask
	err := taskRepo.db.Where("start_date = ? AND religion = ? AND title = ?", date, religion, "Ibadah Minggu").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	var taskCores []entity.ReligionTaskCore
	for _, task := range tasks {
		taskCores = append(taskCores, entity.ReligionTaskModelToTaskCore(task))
	}
	return taskCores, nil
}

// FindAllReligionTask implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllReligionTaskUser(religion string, userId string) ([]entity.ReligionTaskCore, error) {
	var religionTask []model.ReligionTask

	currentTime := time.Now().Format("2006-01-02")

	errData := taskRepo.db.Raw(`
	SELECT * FROM religion_tasks 
	WHERE id NOT IN (
		SELECT task_id FROM user_religion_task_uploads 
		WHERE user_id = ? AND status != 'Ditolak'
	) AND DATE(end_date) >= ?
	AND religion = ?
`, userId, currentTime, religion).Scan(&religionTask).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.ReligionTaskCore, len(religionTask))
	for i, v := range religionTask {
		mapData[i] = entity.ReligionTaskCore{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Point:       v.Point,
			Religion:    v.Religion,
			Start_date:  v.Start_date,
			End_date:    v.End_date,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// FindAllReligionTaskHistory implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllReligionTaskHistory(userId string) ([]entity.UserReligionTaskUploadCore, error) {
	var task []model.UserReligionTaskUpload
	taskRepo.db.Where("user_id=?", userId).Find(&task)

	mapData := make([]entity.UserReligionTaskUploadCore, len(task))
	for i, v := range task {

		userData, _ := taskRepo.userRepository.ReadSpecificUser(userId)
		taskData, _ := taskRepo.FindByIdReligionTask(v.TaskId)

		mapData[i] = entity.UserReligionTaskUploadCore{
			Id:          v.Id,
			UserId:      v.UserId,
			UserName:    userData.Name,
			TaskId:      v.TaskId,
			TaskName:    taskData.Title,
			Image:       v.Image,
			Description: v.Description,
			Type:        v.Type,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// UploadTaskReligion implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UploadTaskReligion(input entity.UserReligionTaskUploadCore, image *multipart.FileHeader) error {
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
	saveDir := "public/images/uploadTaskReligion"
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

	var inputData = model.UserReligionTaskUpload{
		Id:          newUUID,
		UserId:      input.UserId,
		Image:       input.Image,
		TaskId:      input.TaskId,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	errUpload := taskRepo.db.Save(&inputData)
	if errUpload != nil {
		return errUpload.Error
	}

	return nil
}

// FindAllUserReligionTaskUpload implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllUserReligionTaskUpload() ([]entity.UserReligionTaskUploadCore, error) {
	var userTask []model.UserReligionTaskUpload

	errData := taskRepo.db.Find(&userTask).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserReligionTaskUploadCore, len(userTask))
	for i, v := range userTask {
		mapData[i] = entity.UserReligionTaskUploadCore{
			Id:          v.Id,
			TaskId:      v.TaskId,
			UserId:      v.UserId,
			Image:       v.Image,
			Type:        v.Type,
			Description: v.Description,
			Status:      v.Status,
			Message:     v.Message,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// FindSpecificUserReligionTaskUpload implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindSpecificUserReligionTaskUpload(id string) (entity.UserReligionTaskUploadCore, error) {
	var data model.UserReligionTaskUpload

	errData := taskRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserReligionTaskUploadCore{}, errData
	}

	userData, _ := taskRepo.userRepository.ReadSpecificUser(data.UserId)
	taskData, _ := taskRepo.FindByIdReligionTask(data.TaskId)

	userCore := entity.UserReligionTaskUploadCore{
		Id:          data.Id,
		TaskId:      data.TaskId,
		TaskName:    taskData.Title,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Type:        data.Type,
		Message:     data.Message,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return userCore, nil
}

// UpdateReligionTaskStatus implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateReligionTaskStatus(id string, data entity.UserReligionTaskUploadCore) error {
	var pointTask model.ReligionTask
	var userData userModel.Users
	taskData := entity.ReligionTaskUploadCoreToReligionTaskUploadModel(data)

	// get user
	errUser := taskRepo.db.Where("id=?", data.UserId).First(&userData).Error
	if errUser != nil {
		return errUser
	}

	// update status
	tx := taskRepo.db.Where("id=?", data.Id).Updates(taskData)
	if tx.Error != nil {
		return tx.Error
	}

	// get task data
	errData := taskRepo.db.Where("id=?", id).First(&pointTask).Error
	if errData != nil {
		return errData
	}

	if taskData.Status == "Diterima" {
		userPoint, _ := strconv.Atoi(userData.Point)
		userTotalPoint, _ := strconv.Atoi(userData.TotalPoint)

		count := userPoint + pointTask.Point
		countTotal := userTotalPoint + pointTask.Point

		userData.Point = strconv.Itoa(count)
		userData.TotalPoint = strconv.Itoa(countTotal)

		saveUser := user.UserModelToUserCore(userData)

		updateUser := taskRepo.userRepository.UpdatePoint(data.UserId, saveUser)
		if updateUser != nil {
			return updateUser
		}

		//update history
		historyData := user.UserPointCore{
			UserId:   data.UserId,
			Type:     "Religion",
			Point:    pointTask.Point,
			TaskName: pointTask.Title,
		}
		errUserHistory := taskRepo.userRepository.PostUserPointHistory(historyData)
		if errUserHistory != nil {
			return errors.New("failed add user history point")
		}
	}

	return nil
}

// FindAllReligionTaskRequestHistory implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindAllReligionTaskRequestHistory(userId string) ([]entity.UserReligionReqTaskCore, error) {
	var task []model.UserReligionReqTask
	taskRepo.db.Where("user_id=?", userId).Find(&task)

	mapData := make([]entity.UserReligionReqTaskCore, len(task))
	for i, v := range task {

		userData, _ := taskRepo.userRepository.ReadSpecificUser(userId)

		mapData[i] = entity.UserReligionReqTaskCore{
			Id:          v.Id,
			Title:       v.Title,
			UserId:      v.UserId,
			UserName:    userData.Name,
			Image:       v.Image,
			Description: v.Description,
			Status:      v.Status,
			Point:       v.Point,
			Message:     v.Message,
			Type:        v.Type,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// FindSpesificReligionTaskRequest implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) FindSpesificReligionTaskRequest(id string) (entity.UserReligionReqTaskCore, error) {
	var data model.UserReligionReqTask

	errData := taskRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserReligionReqTaskCore{}, errData
	}

	userData, _ := taskRepo.userRepository.ReadSpecificUser(data.UserId)

	userCore := entity.UserReligionReqTaskCore{
		Id:          data.Id,
		Title:       data.Title,
		UserId:      data.UserId,
		UserName:    userData.Name,
		Type:        data.Type,
		Message:     data.Message,
		Point:       data.Point,
		Image:       data.Image,
		Description: data.Description,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return userCore, nil
}

// UploadReligionTaskRequest implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UploadReligionTaskRequest(input entity.UserReligionReqTaskCore, image *multipart.FileHeader) error {
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
	saveDir := "public/images/uploadTaskReligionRequest"
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

	var inputData = model.UserReligionReqTask{
		Id:          newUUID,
		UserId:      input.UserId,
		Title:       input.Title,
		Point:       input.Point,
		Image:       input.Image,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	errUpload := taskRepo.db.Save(&inputData)
	if errUpload != nil {
		return errUpload.Error
	}

	return nil
}

// GetAllUserReligionTaskRequest implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) GetAllUserReligionTaskRequest() ([]entity.UserReligionReqTaskCore, error) {
	var userTask []model.UserReligionReqTask

	errData := taskRepo.db.Find(&userTask).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserReligionReqTaskCore, len(userTask))
	for i, v := range userTask {
		mapData[i] = entity.UserReligionReqTaskCore{
			Id:          v.Id,
			UserId:      v.UserId,
			Title:       v.Title,
			Image:       v.Image,
			Description: v.Description,
			Point:       v.Point,
			Status:      v.Status,
			Message:     v.Message,
			Type:        v.Type,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return mapData, nil
}

// UpdateTaskReligionReqStatus implements entity.TaskDataInterface.
func (taskRepo *TaskRepository) UpdateTaskReligionReqStatus(id string, data entity.UserReligionReqTaskCore) error {
	var pointTask model.UserReligionReqTask
	var userData userModel.Users
	taskData := entity.ReligionTaskReqCoreToReligioinTaskReqModel(data)

	// get user
	errUser := taskRepo.db.Where("id=?", data.UserId).First(&userData).Error
	if errUser != nil {
		return errUser
	}

	// update status
	tx := taskRepo.db.Where("id=?", id).Updates(taskData)
	if tx.Error != nil {
		return tx.Error
	}

	// get task data
	errData := taskRepo.db.Where("id=?", id).First(&pointTask).Error
	if errData != nil {
		return errData
	}

	if taskData.Status == "Diterima" {
		userPoint, _ := strconv.Atoi(userData.Point)
		userTotalPoint, _ := strconv.Atoi(userData.TotalPoint)

		count := userPoint + pointTask.Point
		countTotal := userTotalPoint + pointTask.Point

		userData.Point = strconv.Itoa(count)
		userData.TotalPoint = strconv.Itoa(countTotal)

		saveUser := user.UserModelToUserCore(userData)

		updateUser := taskRepo.userRepository.UpdatePoint(data.UserId, saveUser)
		if updateUser != nil {
			return updateUser
		}

		//update history
		historyData := user.UserPointCore{
			UserId:   data.UserId,
			Type:     "Religion Request",
			Point:    pointTask.Point,
			TaskName: pointTask.Title,
		}
		errUserHistory := taskRepo.userRepository.PostUserPointHistory(historyData)
		if errUserHistory != nil {
			return errors.New("failed add user history point")
		}
	}

	return nil
}
