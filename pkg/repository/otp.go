package repository

import (
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"

	"gorm.io/gorm"
)

type otpRepository struct {
	DB *gorm.DB
}

func NewOtpRepository(DB *gorm.DB) interfaces.OtpRepository {
	return &otpRepository{
		DB: DB,
	}
}

func (ot *otpRepository) FindUserByMobileNumber(phone string) bool {

	var count int
	if err := ot.DB.Raw("select count(*) from users where phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0

}

func (ot *otpRepository) UserDetailsUsingPhone(phone string) (models.UserResponse, error) {

	var usersDetails models.UserResponse
	if err := ot.DB.Raw("select * from users where phone = ?", phone).Scan(&usersDetails).Error; err != nil {
		return models.UserResponse{}, err
	}

	return usersDetails, nil

}
func (ot *otpRepository) UpdateUserPermissionByPhone(phone string) error {
	// Get user details by phone
	var user domain.User
	if err := ot.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return err
	}

	// Check if the user already has the permission
	if user.Permission {
		// User already has the permission, no need to update
		return nil
	}

	// Update the user's permission to true
	if err := ot.DB.Model(&user).Update("permission", true).Error; err != nil {
		return err
	}

	return nil
}
func (ot *otpRepository) ChangePasswordByPhone(phone string, newPassword string) error {
	//a raw SQL query to update the password based on the phone number
	err := ot.DB.Exec("UPDATE users SET password=? WHERE phone=?", newPassword, phone).Error
	if err != nil {
		// If an error occurs during the update, return the error
		return err
	}

	// Return nil if the update is successful
	return nil
}
