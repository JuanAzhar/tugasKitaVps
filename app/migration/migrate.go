package migration

import (
	penalty "tugaskita/features/penalty/model"
	reward "tugaskita/features/reward/model"
	task "tugaskita/features/task/model"
	users "tugaskita/features/user/model"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&users.Users{})
	db.AutoMigrate(&users.UserPoint{})
	db.AutoMigrate(&task.Task{})
	db.AutoMigrate(&task.UserTaskUpload{})
	db.AutoMigrate(&task.UserTaskSubmission{})
	db.AutoMigrate(&reward.Reward{})
	db.AutoMigrate(&reward.UserRewardRequest{})
	db.AutoMigrate(&penalty.Penalty{})
	db.AutoMigrate(&task.ReligionTask{})
	db.AutoMigrate(&task.UserReligionTaskUpload{})
	db.AutoMigrate(&task.UserReligionReqTask{})
}
