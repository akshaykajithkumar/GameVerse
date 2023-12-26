package repository

import (
	"errors"
	"fmt"
	"main/pkg/domain"
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
func (i *userDatabase) ChangePassword(id int, password string) error {

	err := i.DB.Exec("UPDATE users SET password=? WHERE id=?", password, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (i *userDatabase) GetPassword(id int) (string, error) {

	var userPassword string
	err := i.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}

// EditProfile updates the user profile fields in the repository
func (i *userDatabase) EditProfile(id int, name, email, username, phone, bio, url string) error {
	// Create a map to store the fields that need to be updated
	updateFields := make(map[string]interface{})

	// Check if each field is provided and update the map accordingly
	if name != "" {
		updateFields["name"] = name
	}

	if email != "" {
		updateFields["email"] = email
	}

	if username != "" {
		updateFields["username"] = username
	}

	if phone != "" {
		updateFields["phone"] = phone
	}

	if bio != "" {
		updateFields["bio"] = bio
	}
	if url != "" {
		updateFields["url"] = url
	}

	// Update the user's profile fields in the database
	err := i.DB.Model(&domain.User{}).Where("id = ?", id).Updates(updateFields).Error
	if err != nil {
		return err
	}

	return nil
}

// retrieves specific user details by ID from the database
func (i *userDatabase) GetProfileDetailsById(id int) (*domain.User, error) {
	user := &domain.User{}
	err := i.DB.Select("name, email, username, phone, bio,url").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// for storing the reports from userReports
func (i *userDatabase) StoreReport(reporterID, targetID int, reason string) error {
	report := domain.Reports{
		ReporterID: reporterID,
		TargetID:   targetID,
		Reason:     reason,
	}

	// Create a new report record in the database
	if err := i.DB.Create(&report).Error; err != nil {
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
