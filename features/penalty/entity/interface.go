package entity

type PenaltyDataInterface interface {
	CreatePenalty(input PenaltyCore) error
	FindAllPenalty()([]PenaltyCore, error)
	FindSpecificPenalty(id string)(PenaltyCore, error)
	UpdatePenalty(id string, data PenaltyCore) error
	DeletePenalty(id string) error

	FindAllPenaltyHistory(id string)([]PenaltyCore, error)
	GetTotalPenalty(id string)(int,error)
}

type PenaltyUseCaseInterface interface {
	CreatePenalty(input PenaltyCore) error
	FindAllPenalty()([]PenaltyCore, error)
	FindSpecificPenalty(id string)(PenaltyCore, error)
	UpdatePenalty(id string, data PenaltyCore) error
	DeletePenalty(id string) error

	FindAllPenaltyHistory(id string)([]PenaltyCore, error)
	GetTotalPenalty(id string)(int,error)
}