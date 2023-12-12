package repository

import (
	"errors"
	"fmt"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB}
}
func (c *userDatabase) CheckUserAvailability(email string) bool {

	var count int
	query := fmt.Sprintf("select count(*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0
}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	var permission bool
	err := cr.DB.Raw("select permission from users where email = ?", email).Scan(&permission).Error
	if err != nil {
		return false, err
	}
	return permission, nil
}
func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserResponse, error) {

	var user_details models.UserResponse

	err := c.DB.Raw(`
		SELECT *
		FROM users where email = ? and permission = true
		`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserResponse{}, errors.New("error checking user details")
	}

	return user_details, nil
}
func (c *userDatabase) SignUp(user models.UserDetails) (models.UserResponse, error) {

	var userDetails models.UserResponse
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone, username) VALUES (?, ?, ?, ?,?) RETURNING id, name, email, phone", user.Name, user.Email, user.Password, user.Phone, user.Username).Scan(&userDetails).Error

	if err != nil {
		return models.UserResponse{}, err
	}

	return userDetails, nil
}
