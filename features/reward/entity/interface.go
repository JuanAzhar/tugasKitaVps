package entity

import "mime/multipart"

type RewardDataInterface interface {
	CreateReward(input RewardCore, image *multipart.FileHeader) error
	FindAllReward() ([]RewardCore, error)
	FindById(rewardId string) (RewardCore, error)
	UpdateReward(rewardId string, data RewardCore, image *multipart.FileHeader) error
	DeleteReward(rewardId string) error

	UploadRewardRequest(input UserRewardRequestCore) error
	FindAllUploadReward() ([]UserRewardRequestCore, error)
	FindUserRewardById(id string) (UserRewardRequestCore, error)
	FindAllRewardHistory(userId string) ([]UserRewardRequestCore, error)

	UpdateReqRewardStatus(rewardId string, data UserRewardRequestCore) error
}

type RewardUseCaseInterface interface {
	CreateReward(input RewardCore, image *multipart.FileHeader) error
	FindAllReward() ([]RewardCore, error)
	FindById(rewardId string) (RewardCore, error)
	UpdateReward(rewardId string, data RewardCore, image *multipart.FileHeader) error
	DeleteReward(rewardId string) error

	UploadRewardRequest(input UserRewardRequestCore) error
	FindAllUploadReward() ([]UserRewardRequestCore, error)
	FindUserRewardById(id string) (UserRewardRequestCore, error)
	FindAllRewardHistory(userId string) ([]UserRewardRequestCore, error)

	UpdateReqRewardStatus(rewardId string, data UserRewardRequestCore) error
}
