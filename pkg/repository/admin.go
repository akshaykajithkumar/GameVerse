package repository

import (
	"fmt"
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"
	"strconv"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}
func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {

	var adminCompareDetails domain.Admin
	if err := ad.DB.Raw("select * from admins where email = ? ", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Admin{}, err
	}
	return adminCompareDetails, nil
}

func (ad *adminRepository) GetUsers(page int, limit int) ([]models.UserDetailsAtAdmin, error) {
	// pagination purpose -
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	var userDetails []models.UserDetailsAtAdmin

	if err := ad.DB.Raw("select id,name,email,phone,permission from users limit ? offset ?", limit, offset).Scan(&userDetails).Error; err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

func (ad *adminRepository) GetUserByID(id string) (domain.User, error) {

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, err
	}

	query := fmt.Sprintf("select * from users where id = '%d'", user_id)
	var userDetails domain.User

	if err := ad.DB.Raw(query).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}

	return userDetails, nil
}

// function which will both block and unblock a user
func (ad *adminRepository) UpdateBlockUserByID(user domain.User) error {

	err := ad.DB.Exec("update users set permission = ? where id = ?", user.Permission, user.ID).Error
	if err != nil {
		return err
	}

	return nil

}
func (i *adminRepository) GetReports(page, limit int) ([]domain.Reports, error) {
	var reports []domain.Reports
	offset := (page - 1) * limit

	if err := i.DB.Offset(offset).Limit(limit).Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}
func (ar *adminRepository) AddSubscriptionPlan(name string, duration int, price float64) error {
	// Build the raw SQL query
	query := "INSERT INTO subscription_plans (name, duration, price) VALUES (?, ?, ?)"

	// Execute the raw SQL query
	result := ar.DB.Exec(query, name, duration, price)
	if result.Error != nil {
		// Handle any error during the execution, you can log or perform additional actions as needed
		return result.Error
	}

	// Return nil if the subscription plan was inserted successfully
	return nil
}
func (ar *adminRepository) DeleteSubscriptionPlan(planID int) error {
	// Build the raw SQL query to delete the subscription plan by ID
	query := "DELETE FROM subscription_plans WHERE id = ?"

	// Execute the raw SQL query
	result := ar.DB.Exec(query, planID)
	if result.Error != nil {
		// Handle any error during the execution, you can log or perform additional actions as needed
		return result.Error
	}

	// Return nil if the subscription plan was deleted successfully
	return nil
}
func (ar *adminRepository) GetSubscriptionPlans() ([]domain.SubscriptionPlan, error) {
	var plans []domain.SubscriptionPlan

	// Fetch the list of subscription plans from the database
	err := ar.DB.Find(&plans).Error
	if err != nil {
		// Handle any error during the fetch operation, you can log or perform additional actions as needed
		return nil, err
	}

	// Return the list of subscription plans
	return plans, nil
}
