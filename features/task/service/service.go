package service

import (
	"errors"
	"mime/multipart"
	"time"
	"tugaskita/features/task/entity"
)

type taskService struct {
	TaskRepo entity.TaskDataInterface
}

func NewTaskService(taskRepo entity.TaskDataInterface) entity.TaskUseCaseInterface {
	return &taskService{
		TaskRepo: taskRepo,
	}
}

// CreateTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) CreateTask(data entity.TaskCore) error {
	layout := "2006-01-02"
	currentTime := time.Now().Truncate(24 * time.Hour)

	if data.Title == "" || data.Description == "" {
		return errors.New("title and description can't empty")
	}

	start, errStart := time.Parse(layout, data.Start_date)
	if errStart != nil {
		return errors.New("start date must be in 'yyyy-mm-dd' format")
	}
	if start.Before(currentTime) {
		return errors.New("please choose at least today")
	}

	end, errEnd := time.Parse(layout, data.End_date)
	if errEnd != nil {
		return errors.New("end date must be in 'yyyy-mm-dd' format")
	}

	if end.Before(start) {
		return errors.New("end date must be after start date")
	}

	if end.Equal(start) {
		return errors.New("end date must be different from start date")
	}

	if data.Point <= 0 {
		return errors.New("point must be more than 0")
	}

	err := taskUC.TaskRepo.CreateTask(data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask implements entity.TaskCoreUseCaseInterface.
func (taskUC *taskService) DeleteTask(taskId string) error {
	if taskId == "" {
		return errors.New("insert task id")
	}

	_, err := taskUC.TaskRepo.FindById(taskId)
	if err != nil {
		return errors.New("task not found")
	}

	errDelete := taskUC.TaskRepo.DeleteTask(taskId)
	if errDelete != nil {
		return errors.New("can't delete task")
	}

	return nil
}

// FindAllMission implements entity.TaskCoreUseCaseInterface.
func (taskUC *taskService) FindAllTask() ([]entity.TaskCore, error) {
	data, err := taskUC.TaskRepo.FindAllTask()
	if err != nil {
		return nil, errors.New("error get data")
	}

	return data, nil
}

// FindById implements entity.TaskCoreUseCaseInterface.
func (taskUC *taskService) FindById(taskId string) (entity.TaskCore, error) {
	if taskId == "" {
		return entity.TaskCore{}, errors.New("task ID is required")
	}

	task, err := taskUC.TaskRepo.FindById(taskId)
	if err != nil {
		return entity.TaskCore{}, err
	}

	return task, nil
}

// UpdateTask implements entity.TaskCoreUseCaseInterface.
func (taskUC *taskService) UpdateTask(taskId string, data entity.TaskCore) error {
	layout := "2006-01-02"
	currentTime := time.Now().Truncate(24 * time.Hour)

	if data.Title == "" || data.Description == "" {
		return errors.New("title and description can't empty")
	}

	start, errStart := time.Parse(layout, data.Start_date)
	if errStart != nil {
		return errors.New("start date must be in 'yyyy-mm-dd' format")
	}
	if start.Before(currentTime) {
		return errors.New("please choose at least today")
	}

	end, errEnd := time.Parse(layout, data.End_date)
	if errEnd != nil {
		return errors.New("end date must be in 'yyyy-mm-dd' format")
	}

	if end.Before(start) {
		return errors.New("end date must be after start date")
	}

	if end.Equal(start) {
		return errors.New("end date must be different from start date")
	}

	if data.Point <= 0 {
		return errors.New("point must be more than 0")
	}

	err := taskUC.TaskRepo.UpdateTask(taskId, data)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTaskStatus implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UpdateTaskStatus(taskId string, data entity.UserTaskUploadCore) error {
	if data.Status == "" {
		return errors.New("status can't be empty")
	}

	taskData, errData := taskUC.TaskRepo.FindUserTaskById(taskId)
	if errData != nil {
		return errors.New("task not found")
	}

	if data.Status == taskData.Status {
		return errors.New("you already updated this task to " + data.Status)
	}

	err := taskUC.TaskRepo.UpdateTaskStatus(taskId, data)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTaskReqStatus implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UpdateTaskReqStatus(id string, data entity.UserTaskSubmissionCore) error {
	if data.Status == "" {
		return errors.New("status can't be empty")
	}

	taskData, errData := taskUC.TaskRepo.FindUserTaskReqById(id)
	if errData != nil {
		return errors.New("request task not found")
	}

	if data.Status == taskData.Status {
		return errors.New("you already updated this task to " + data.Status)
	}

	err := taskUC.TaskRepo.UpdateTaskReqStatus(id, data)
	if err != nil {
		return err
	}

	return nil
}

// FindAllClaimedTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllClaimedTask(userId string) ([]entity.UserTaskUploadCore, error) {
	data, err := taskUC.TaskRepo.FindAllClaimedTask(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindAllRequestTaskHistory implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllRequestTaskHistory(userId string) ([]entity.UserTaskSubmissionCore, error) {
	data, err := taskUC.TaskRepo.FindAllRequestTaskHistory(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindTasksNotClaimedByUser implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindTasksNotClaimedByUser(userId string) ([]entity.TaskCore, error) {
	data, err := taskUC.TaskRepo.FindTasksNotClaimedByUser(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UploadTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UploadTask(data entity.UserTaskUploadCore, image *multipart.FileHeader) error {
	_, errTask := taskUC.TaskRepo.FindById(data.TaskId)
	if errTask != nil {
		return errors.New("task not found")
	}

	if data.Description == "" {
		return errors.New("description can't empty")
	}

	if image != nil && image.Size > 10*1024*1024 {
		return errors.New("image file size should be less than 10 MB")
	}

	err := taskUC.TaskRepo.UploadTask(data, image)
	if err != nil {
		return errors.New("failed upload task")
	}

	return nil
}

// FindUserTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllUserTask() ([]entity.UserTaskUploadCore, error) {
	userTask, err := taskUC.TaskRepo.FindAllUserTask()
	if err != nil {
		return nil, errors.New("error get user task")
	}

	return userTask, nil
}

// FindUserTaskById implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindUserTaskById(id string) (entity.UserTaskUploadCore, error) {
	task, err := taskUC.TaskRepo.FindUserTaskById(id)
	if err != nil {
		return entity.UserTaskUploadCore{}, err
	}

	return task, nil
}

// FindUserTaskReqById implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindUserTaskReqById(id string) (entity.UserTaskSubmissionCore, error) {
	task, err := taskUC.TaskRepo.FindUserTaskReqById(id)
	if err != nil {
		return entity.UserTaskSubmissionCore{}, err
	}

	return task, nil
}

// UploadTaskRequest implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UploadTaskRequest(input entity.UserTaskSubmissionCore, image *multipart.FileHeader) error {
	if input.Description == "" || input.Title == "" {
		return errors.New("description and title can't empty")
	}

	if input.Point <= 0 {
		return errors.New("point can't less then 0")
	}

	if image != nil && image.Size > 10*1024*1024 {
		return errors.New("image file size should be less than 10 MB")
	}

	err := taskUC.TaskRepo.UploadTaskRequest(input, image)
	if err != nil {
		return errors.New("failed upload request task")
	}

	return nil
}

// FindAllRequestTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllRequestTask() ([]entity.UserTaskSubmissionCore, error) {
	userTask, err := taskUC.TaskRepo.FindAllRequestTask()
	if err != nil {
		return nil, errors.New("error get user request task")
	}

	return userTask, nil
}

// CountUserClearTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) CountUserClearTask(id string) (int, error) {
	countTask, err := taskUC.TaskRepo.CountUserClearTask(id)
	if err != nil {
		return 0, errors.New("error count user clear task")
	}

	return countTask, nil
}

// CreateTaskReligion implements entity.TaskUseCaseInterface.
func (taskUC *taskService) CreateTaskReligion(input entity.ReligionTaskCore) error {
	layout := "2006-01-02"
	currentTime := time.Now().Truncate(24 * time.Hour)

	if input.Point < 0 {
		return errors.New("point must be more than 0")
	}

	if input.Religion == "" {
		return errors.New("religion can't empty")
	}

	// Kondisi khusus untuk agama Islam
	if input.Religion == "Islam" && input.Title == "" {
		//cek apakah hari ini sudah upload atau belum
		existingTasks, err := taskUC.TaskRepo.FindTaskByDateAndReligion(currentTime.Format(layout), "Islam")
		if err != nil {
			return err
		}

		if len(existingTasks) > 0 {
			return errors.New("shalat 5 waktu sudah dibuat untuk hari ini")
		}

		dayData := time.Now().Weekday().String()
		println(dayData)

		if dayData == "Friday" {
			//post shalat 5 waktu
			prayers := []string{"Subuh", "Jum'at", "Ashar", "Maghrib", "Isya"}
			for _, prayer := range prayers {
				task := entity.ReligionTaskCore{
					Title:       prayer,
					Point:       250,
					Religion:    input.Religion,
					Start_date:  currentTime.Format(layout),
					End_date:    currentTime.Format(layout),
					Description: "Tugas Shalat " + prayer,
				}
				err := taskUC.TaskRepo.CreateTaskReligion(task)
				if err != nil {
					return err
				}
			}
		}

		//post shalat 5 waktu
		prayers := []string{"Subuh", "Dzuhur", "Ashar", "Maghrib", "Isya"}
		for _, prayer := range prayers {
			task := entity.ReligionTaskCore{
				Title:       prayer,
				Point:       250,
				Religion:    input.Religion,
				Start_date:  currentTime.Format(layout),
				End_date:    currentTime.Format(layout),
				Description: "Tugas Shalat " + prayer,
			}
			err := taskUC.TaskRepo.CreateTaskReligion(task)
			if err != nil {
				return err
			}
		}
	} else if input.Religion == "Kristen" && input.Title == "" {
		currentWeekday := time.Now().Weekday()

		// Hitung jarak hari ke hari Minggu berikutnya
		daysUntilSunday := (7 - int(currentWeekday))

		// Jika hari ini bukan hari Minggu, tambahkan jarak hari ke tanggal saat ini
		if daysUntilSunday != 0 {
			currentTime = currentTime.AddDate(0, 0, daysUntilSunday)
		}

		//cek apakah hari ini sudah upload atau belum
		existingTasks, errCheck := taskUC.TaskRepo.FindTaskByDateAndReligionNon(currentTime.Format(layout), "Kristen")
		if errCheck != nil {
			return errCheck
		}

		if len(existingTasks) > 0 {
			return errors.New("ibadah minggu telah ditambahkan pada minggu ini")
		}

		task := entity.ReligionTaskCore{
			Title:       "Ibadah Minggu",
			Point:       300,
			Religion:    input.Religion,
			Start_date:  currentTime.Format(layout),
			End_date:    currentTime.Format(layout),
			Description: "Laksanakan ibadah minggu ke gereja",
		}
		err := taskUC.TaskRepo.CreateTaskReligion(task)
		if err != nil {
			return err
		}
	} else if input.Religion == "Katolik" && input.Title == "" {
		currentWeekday := time.Now().Weekday()

		// Hitung jarak hari ke hari Minggu berikutnya
		daysUntilSunday := (7 - int(currentWeekday))

		// Jika hari ini bukan hari Minggu, tambahkan jarak hari ke tanggal saat ini
		if daysUntilSunday != 0 {
			currentTime = currentTime.AddDate(0, 0, daysUntilSunday)
		}

		//cek apakah hari ini sudah upload atau belum
		existingTasks, errCheck := taskUC.TaskRepo.FindTaskByDateAndReligionNon(currentTime.Format(layout), "Katolik")
		if errCheck != nil {
			return errCheck
		}

		if len(existingTasks) > 0 {
			return errors.New("ibadah minggu telah ditambahkan pada minggu ini")
		}

		task := entity.ReligionTaskCore{
			Title:       "Ibadah Minggu",
			Point:       300,
			Religion:    input.Religion,
			Start_date:  currentTime.Format(layout),
			End_date:    currentTime.Format(layout),
			Description: "Laksanakan ibadah minggu ke gereja",
		}
		err := taskUC.TaskRepo.CreateTaskReligion(task)
		if err != nil {
			return err
		}
	} else {

		if input.Religion == "" || input.Title == "" {
			return errors.New("religion and title can't empty")
		}

		start, errStart := time.Parse(layout, input.Start_date)

		if errStart != nil {
			return errors.New("start date must be in 'yyyy-mm-dd' format")
		}
		if start.Before(currentTime) {
			return errors.New("please choose at least today")
		}

		end, errEnd := time.Parse(layout, input.End_date)
		if errEnd != nil {
			return errors.New("end date must be in 'yyyy-mm-dd' format")
		}

		if end.Before(start) {
			return errors.New("end date must be after start date")
		}

		if end.Equal(start) {
			return errors.New("end date must be different from start date")
		}

		err := taskUC.TaskRepo.CreateTaskReligion(input)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteTaskReligion implements entity.TaskUseCaseInterface.
func (taskUC *taskService) DeleteTaskReligion(taskId string) error {
	if taskId == "" {
		return errors.New("insert task id")
	}

	_, err := taskUC.TaskRepo.FindByIdReligionTask(taskId)
	if err != nil {
		return errors.New("task not found")
	}

	errDelete := taskUC.TaskRepo.DeleteTaskReligion(taskId)
	if errDelete != nil {
		return errors.New("can't delete task")
	}

	return nil
}

// FindAllTaskReligion implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllTaskReligion() ([]entity.ReligionTaskCore, error) {
	data, err := taskUC.TaskRepo.FindAllTaskReligion()
	if err != nil {
		return nil, errors.New("error get data")
	}

	return data, nil
}

// FindByIdReligionTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindByIdReligionTask(taskId string) (entity.ReligionTaskCore, error) {
	if taskId == "" {
		return entity.ReligionTaskCore{}, errors.New("task ID is required")
	}

	task, err := taskUC.TaskRepo.FindByIdReligionTask(taskId)
	if err != nil {
		return entity.ReligionTaskCore{}, err
	}

	return task, nil
}

// UpdateTaskReligion implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UpdateTaskReligion(taskId string, data entity.ReligionTaskCore) error {
	if data.Point < 0 {
		return errors.New("point must be more than 0")
	}

	if data.Religion == "" || data.Title == "" {
		return errors.New("religion and title can't empty")
	}

	err := taskUC.TaskRepo.UpdateTaskReligion(taskId, data)
	if err != nil {
		return err
	}

	return nil
}

// FindAllReligionTask implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllReligionTaskUser(religion string, userId string) ([]entity.ReligionTaskCore, error) {
	religionTask, err := taskUC.TaskRepo.FindAllReligionTaskUser(religion, userId)
	if err != nil {
		return nil, errors.New("error get religion task")
	}

	return religionTask, nil
}

// FindAllReligionTaskHistory implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllReligionTaskHistory(userId string) ([]entity.UserReligionTaskUploadCore, error) {
	data, err := taskUC.TaskRepo.FindAllReligionTaskHistory(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UploadTaskReligion implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UploadTaskReligion(input entity.UserReligionTaskUploadCore, image *multipart.FileHeader) error {
	_, errTask := taskUC.TaskRepo.FindByIdReligionTask(input.TaskId)
	if errTask != nil {
		return errors.New("religion task not found")
	}

	if input.Description == "" {
		return errors.New("description can't empty")
	}

	if image != nil && image.Size > 10*1024*1024 {
		return errors.New("image file size should be less than 10 MB")
	}

	err := taskUC.TaskRepo.UploadTaskReligion(input, image)
	if err != nil {
		return errors.New("failed upload religion task")
	}

	return nil
}

// FindAllUserReligionTaskUpload implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllUserReligionTaskUpload() ([]entity.UserReligionTaskUploadCore, error) {
	userTask, err := taskUC.TaskRepo.FindAllUserReligionTaskUpload()
	if err != nil {
		return nil, errors.New("error get user religion task data")
	}

	return userTask, nil
}

// FindSpecificUserReligionTaskUpload implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindSpecificUserReligionTaskUpload(id string) (entity.UserReligionTaskUploadCore, error) {
	task, err := taskUC.TaskRepo.FindSpecificUserReligionTaskUpload(id)
	if err != nil {
		return entity.UserReligionTaskUploadCore{}, err
	}

	return task, nil
}

// UpdateReligionTaskStatus implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UpdateReligionTaskStatus(id string, data entity.UserReligionTaskUploadCore) error {
	if data.Status == "" {
		return errors.New("status can't be empty")
	}

	taskData, errData := taskUC.TaskRepo.FindSpecificUserReligionTaskUpload(data.Id.String())
	if errData != nil {
		return errors.New("religion task not found")
	}

	if data.Status == taskData.Status {
		return errors.New("you already updated this task to " + data.Status)
	}

	err := taskUC.TaskRepo.UpdateReligionTaskStatus(id, data)
	if err != nil {
		return err
	}

	return nil
}

// FindAllReligionTaskRequestHistory implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindAllReligionTaskRequestHistory(userId string) ([]entity.UserReligionReqTaskCore, error) {
	data, err := taskUC.TaskRepo.FindAllReligionTaskRequestHistory(userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FindSpesificReligionTaskRequest implements entity.TaskUseCaseInterface.
func (taskUC *taskService) FindSpesificReligionTaskRequest(id string) (entity.UserReligionReqTaskCore, error) {
	task, err := taskUC.TaskRepo.FindSpesificReligionTaskRequest(id)
	if err != nil {
		return entity.UserReligionReqTaskCore{}, err
	}

	return task, nil
}

// UploadReligionTaskRequest implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UploadReligionTaskRequest(input entity.UserReligionReqTaskCore, image *multipart.FileHeader) error {
	if input.Description == "" || input.Title == "" {
		return errors.New("description and title can't empty")
	}

	if input.Point <= 0 {
		return errors.New("point can't less then 0")
	}

	if image != nil && image.Size > 10*1024*1024 {
		return errors.New("image file size should be less than 10 MB")
	}

	err := taskUC.TaskRepo.UploadReligionTaskRequest(input, image)
	if err != nil {
		return errors.New("failed upload religion request task")
	}

	return nil
}

// GetAllUserReligionTaskRequest implements entity.TaskUseCaseInterface.
func (taskUC *taskService) GetAllUserReligionTaskRequest() ([]entity.UserReligionReqTaskCore, error) {
	userTask, err := taskUC.TaskRepo.GetAllUserReligionTaskRequest()
	if err != nil {
		return nil, errors.New("error get user religion request task")
	}

	return userTask, nil
}

// UpdateTaskReligionReqStatus implements entity.TaskUseCaseInterface.
func (taskUC *taskService) UpdateTaskReligionReqStatus(id string, data entity.UserReligionReqTaskCore) error {
	if data.Status == "" {
		return errors.New("status can't be empty")
	}

	taskData, errData := taskUC.TaskRepo.FindSpesificReligionTaskRequest(id)
	if errData != nil {
		return errors.New("religion request task not found")
	}

	if data.Status == taskData.Status {
		return errors.New("you already updated this task to " + data.Status)
	}

	err := taskUC.TaskRepo.UpdateTaskReligionReqStatus(id, data)
	if err != nil {
		return err
	}

	return nil
}