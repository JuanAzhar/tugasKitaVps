package dto

type RewardRequest struct {
	Name  string `json:"name" form:"name"`
	Stock int    `json:"stock" form:"stock"`
	Price int    `json:"price" form:"price"`
	Image string `json:"image" form:"image"`
}

type RewardReqRequest struct {
	RewardId string `json:"reward_id"`
	Amount   int    `json:"amount"`
}

type RewardReqUpdateRequest struct {
	RewardId string `json:"reward_id"`
	UserId   string `json:"user_id"`
	Status   string `json:"status"`
}
