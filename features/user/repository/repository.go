package repository

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"tugaskita/features/user/entity"
	"tugaskita/features/user/model"
	bcrypt "tugaskita/utils/bcrypt"
	utils "tugaskita/utils/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) entity.UserDataInterface {
	return &userRepository{
		db: db,
	}
}

// DeleteUser implements entity.UserDataInterface.
func (userRepo *userRepository) DeleteUser(id string) (err error) {
	var chekcId model.Users

	errData := userRepo.db.Where("id = ?", id).Delete(&chekcId)
	if errData != nil {
		return errData.Error
	}

	return nil

}

// Login implements entity.UserDataInterface.
func (userRepo *userRepository) Login(email string, password string) (entity.UserCore, string, error) {
	var data model.Users

	bcrypt.CheckPasswordHash(data.Password, password)

	tx := userRepo.db.Where("email=?", email).First(&data)
	if tx.Error != nil {
		return entity.UserCore{}, "", tx.Error
	}

	var token string

	if tx.RowsAffected > 0 {
		var errToken error
		token, errToken = utils.CreateToken(data.ID, data.Role, data.Religion)
		if errToken != nil {
			return entity.UserCore{}, "", errToken
		}
	}

	var resp = entity.UserCore{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
	}

	return resp, token, nil
}

// ReadSpecificUser implements entity.UserDataInterface.
func (userRepo *userRepository) ReadSpecificUser(id string) (user entity.UserCore, err error) {
	var data model.Users
	errData := userRepo.db.Where("id=?", id).First(&data).Error
	if errData != nil {
		return entity.UserCore{}, errData
	}

	userCore := entity.UserCore{
		ID:         data.ID,
		Name:       data.Name,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Email:      data.Email,
		Image:      data.Image,
		Religion:   data.Religion,
		Point:      data.Point,
		TotalPoint: data.TotalPoint,
		Role:       data.Role,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}

	return userCore, nil
}

func (userRepo *userRepository) Register(data entity.UserCore, image *multipart.FileHeader) (row int, err error) {
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return 0, UUIDerr
	}

	hashPassword, err := bcrypt.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	file, err := image.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Define the directory where you want to save the image
	saveDir := "public/images/user"
	os.MkdirAll(saveDir, os.ModePerm)

	// Define the file path
	filePath := filepath.Join(saveDir, image.Filename)

	filePath = strings.ReplaceAll(filePath, "\\", "/")

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	// Copy the uploaded file data to the destination file
	if _, err = io.Copy(dst, file); err != nil {
		return 0, err
	}

	// Set the image URL to the file path
	data.Image = filePath

	var input = model.Users{
		ID:         newUUID.String(),
		Name:       data.Name,
		Image:      data.Image,
		Address:    data.Address,
		School:     data.School,
		Class:      data.Class,
		Email:      data.Email,
		Religion:   data.Religion,
		Password:   hashPassword,
		Point:      "0",
		TotalPoint: "0",
		Role:       "user",
	}

	println(input.Address)
	println(input.School)
	println(input.Class)

	erruser := userRepo.db.Save(&input)
	if erruser.Error != nil {
		return 0, erruser.Error
	}

	return 1, nil
}

// ReadAllUser implements entity.UserDataInterface.
func (userRepo *userRepository) ReadAllUser() ([]entity.UserCore, error) {
	var dataUser []model.Users

	errData := userRepo.db.Where("role = ?", "user").Find(&dataUser).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserCore, len(dataUser))
	for i, value := range dataUser {
		mapData[i] = entity.UserCore{
			ID:         value.ID,
			Name:       value.Name,
			Email:      value.Email,
			Address:    value.Address,
			School:     value.School,
			Class:      value.Class,
			Image:      value.Image,
			Role:       value.Role,
			Religion:   value.Religion,
			Point:      value.Point,
			TotalPoint: value.TotalPoint,
		}
	}
	return mapData, nil
}

// // UpdateSiswa implements entity.UserDataInterface.
// func (userRepo *userRepository) UpdateSiswa(id string, data entity.UserCore, image *multipart.FileHeader) error {
// 	dataUser := entity.UserCoreToUserModel(data)

// 	// Jika gambar diunggah, lakukan upload ke Cloudinary
// 	if image != nil {
// 		file, err := image.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		imageURL, err := cloudinary.UploadToCloudinary(file, image.Filename)
// 		if err != nil {
// 			return err
// 		}

// 		dataUser.Image = imageURL
// 	}

// 	// Hash password jika tidak kosong
// 	if data.Password != "" {
// 		hashPassword, err := bcrypt.HashPassword(data.Password)
// 		if err != nil {
// 			return err
// 		}
// 		dataUser.Password = hashPassword
// 	}

// 	// Update data pengguna
// 	tx := userRepo.db.Where("id = ?", id).Updates(&dataUser)
// 	if tx.Error != nil {
// 		return tx.Error
// 	}

// 	if tx.RowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	return nil
// }

func (userRepo *userRepository) UpdateSiswa(id string, data entity.UserCore, image *multipart.FileHeader) error {
	dataUser := entity.UserCoreToUserModel(data)

	// If an image is uploaded, save it to a local folder
	if image != nil {
		file, err := image.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// Define the directory where you want to save the image
		saveDir := "public/images/user"
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

		// Update the dataUser.Image field with the local file path
		dataUser.Image = filePath
	}

	// Hash the password if it's not empty
	if data.Password != "" {
		hashPassword, err := bcrypt.HashPassword(data.Password)
		if err != nil {
			return err
		}
		dataUser.Password = hashPassword
	}

	// Update the user's data in the database
	tx := userRepo.db.Where("id = ?", id).Updates(&dataUser)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdatePoint implements entity.UserDataInterface.
func (userRepo *userRepository) UpdatePoint(id string, data entity.UserCore) error {
	userData := entity.UserCoreToUserModel(data)

	tx := userRepo.db.Where("id = ?", id).Updates(&userData)
	if tx != nil {
		return tx.Error
	}

	return nil
}

// GetRankUser implements entity.UserDataInterface.
func (userRepo *userRepository) GetRankUser() ([]entity.UserCore, error) {
	var dataUser []model.Users

	errData := userRepo.db.Where("role = ?", "user").Order("point desc").Find(&dataUser).Error
	if errData != nil {
		return nil, errData
	}

	mapData := make([]entity.UserCore, len(dataUser))
	for i, value := range dataUser {
		mapData[i] = entity.UserCore{
			ID:         value.ID,
			Name:       value.Name,
			Email:      value.Email,
			Role:       value.Role,
			Point:      value.Point,
			TotalPoint: value.TotalPoint,
		}
	}
	return mapData, nil
}

// ChangePassword implements entity.UserDataInterface.
func (userRepo *userRepository) ChangePassword(id string, data entity.UserCore) error {
	password := entity.UserCoreToUserModel(data)

	hashPassword, err := bcrypt.HashPassword(data.Password)
	if err != nil {
		return err
	}

	password.Password = hashPassword

	tx := userRepo.db.Where("id = ?", id).Updates(&password)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// AnnualResetPoint implements entity.UserDataInterface.
func (userRepo *userRepository) AnnualResetPoint() error {
	err := userRepo.db.Exec("UPDATE users SET total_point = 0").Error
	if err != nil {
		return err
	}
	return nil
}

// MonthlyResetPoint implements entity.UserDataInterface.
func (userRepo *userRepository) MonthlyResetPoint() error {
	err := userRepo.db.Exec("UPDATE users SET point = 0").Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllUserPointHistory implements entity.UserDataInterface.
func (userRepo *userRepository) GetAllUserPointHistory() ([]entity.UserPointCore, error) {
	var userPoint []model.UserPoint
	userRepo.db.Find(&userPoint)

	dataUser := entity.ListUserPointModelToListUserPointCore(userPoint)
	return dataUser, nil
}

// GetSpecificUserPointHistory implements entity.UserDataInterface.
func (userRepo *userRepository) GetSpecificUserPointHistory(id string) (entity.UserPointCore, error) {
	dataUser := model.UserPoint{}

	tx := userRepo.db.Where("id = ? ", id).First(&dataUser)
	if tx.Error != nil {
		return entity.UserPointCore{}, tx.Error
	}

	dataResponse := entity.UserPointModelToUserPointCore(dataUser)
	return dataResponse, nil
}

// GetUserPointHistory implements entity.UserDataInterface.
func (userRepo *userRepository) GetUserPointHistory(id string) ([]entity.UserPointCore, error) {
	var userPoint []model.UserPoint
	userRepo.db.Where("user_id=?", id).Find(&userPoint)

	datauserPoint := entity.ListUserPointModelToListUserPointCore(userPoint)
	return datauserPoint, nil
}

// PostUserPointHistory implements entity.UserDataInterface.
func (userRepo *userRepository) PostUserPointHistory(data entity.UserPointCore) error {
	println("masuk kesini")
	newUUID, UUIDerr := uuid.NewRandom()
	if UUIDerr != nil {
		return UUIDerr
	}

	data.Id = newUUID.String()
	dataUser := entity.UserPointCoreToUserPointModel(data)

	tx := userRepo.db.Create(&dataUser)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
