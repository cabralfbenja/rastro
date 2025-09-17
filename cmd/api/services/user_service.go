package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"rastro/cmd/api/requests"
	"rastro/common"
	"rastro/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) Register(userRequest *requests.RegisterUserRequest) (*models.UserModel, error) {
	// Hash password
	hashedPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("User registration failed")
	}
	createdUser := models.UserModel{
		Username:     userRequest.Username,
		Email:        userRequest.Email,
		FirstName:    &userRequest.FirstName,
		LastName:     &userRequest.LastName,
		PasswordHash: hashedPassword,
	}

	result := us.db.Create(&createdUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("User registration failed")
	}

	return &createdUser, nil
}

func (us *UserService) GetUserByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel
	result := us.db.First(&user, "username = ?", username)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (us *UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := us.db.First(&user, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (us *UserService) ChangeUserPassword(newPassword string, user models.UserModel) error {
	hashedNewPassword, err := common.HashPassword(newPassword)
	if err != nil {
		return errors.New("Failed to update password")
	}

	user.PasswordHash = hashedNewPassword
	result := us.db.Model(&user).Update("password_hash", hashedNewPassword)
	if result.Error != nil {
		return errors.New("Failed to update password")
	}
	return nil

}
