package entity

import "mime/multipart"

type UserDataInterface interface {
	Register(data UserCore, image *multipart.FileHeader) (row int, err error)
	UpdateSiswa(id string, data UserCore, image *multipart.FileHeader) error
	Login(email, password string) (UserCore, string, error)
	ReadAllUser() ([]UserCore, error)
	ReadSpecificUser(id string) (user UserCore, err error)
	DeleteUser(id string) (err error)
	UpdatePoint(id string, data UserCore) error

	GetRankUser() ([]UserCore, error)
	ChangePassword(id string, data UserCore) error

	MonthlyResetPoint()(error)
	AnnualResetPoint()(error)

	PostUserPointHistory(data UserPointCore) error
	GetAllUserPointHistory()([]UserPointCore, error)
	GetSpecificUserPointHistory(id string)(UserPointCore, error)
	GetUserPointHistory(id string)([]UserPointCore, error)
}

type UserUseCaseInterface interface {
	Register(data UserCore, image *multipart.FileHeader) (row int, err error)
	UpdateSiswa(id string, data UserCore, image *multipart.FileHeader) error
	Login(email, password string) (UserCore, string, error)
	ReadAllUser() ([]UserCore, error)
	ReadSpecificUser(id string) (user UserCore, err error)
	DeleteUser(id string) (err error)

	GetRankUser() ([]UserCore, error)
	ChangePassword(id string, data UserCore) error

	MonthlyResetPoint()(error)
	AnnualResetPoint()(error)
	
	PostUserPointHistory(data UserPointCore) error
	GetAllUserPointHistory()([]UserPointCore, error)
	GetSpecificUserPointHistory(id string)(UserPointCore, error)
	GetUserPointHistory(id string)([]UserPointCore, error)
}
