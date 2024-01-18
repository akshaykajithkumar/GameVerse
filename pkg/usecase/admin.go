package usecase

import (
	"errors"
	"main/pkg/domain"
	"main/pkg/helper"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"

	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository interfaces.AdminRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository) services.AdminUseCase {
	return &adminUseCase{
		adminRepository: repo,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (models.TokenAdmin, error) {

	// getting details of the admin based on the email provided
	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return models.TokenAdmin{}, err
	}

	// compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return models.TokenAdmin{}, err
	}

	accessTokenString, refreshTokenString, err := helper.GenerateTokensAdmin(adminCompareDetails)

	if err != nil {
		return models.TokenAdmin{}, err
	}

	return models.TokenAdmin{
		Username:     adminCompareDetails.Username,
		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	}, nil

}
func (ad *adminUseCase) GetUsers(page int, limit int) ([]models.UserDetailsAtAdmin, error) {

	userDetails, err := ad.adminRepository.GetUsers(page, limit)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

func (ad *adminUseCase) ToggleBlockUser(id string) error {
	// Retrieve user by ID
	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	// Toggle the Permission field
	user.Permission = !user.Permission

	// Update user in the repository
	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *adminUseCase) GetReports(page, limit int) ([]domain.Reports, error) {
	reports, err := u.adminRepository.GetReports(page, limit)
	if err != nil {
		return nil, err
	}
	return reports, nil
}
func (u *adminUseCase) AddSubscriptionPlan(name string, duration int, price float64) error {
	// Validate the input parameters
	if duration <= 10 || price <= 10 {
		// Return an error if any of the values are not above 10
		return errors.New("duration and price must be above 10")
	}

	// Perform any additional business logic or validation specific to adding a subscription plan

	// Call the repository function to add the subscription plan
	err := u.adminRepository.AddSubscriptionPlan(name, duration, price)
	if err != nil {
		// Handle any error from the repository, you can log or perform additional actions as needed
		return err
	}

	// Return nil if the subscription plan was added successfully
	return nil
}
func (u *adminUseCase) DeleteSubscriptionPlan(planID int) error {
	// Perform any business logic or validation specific to deleting a subscription plan

	// Call the repository function to delete the subscription plan
	err := u.adminRepository.DeleteSubscriptionPlan(planID)
	if err != nil {
		// Handle any error from the repository, you can log or perform additional actions as needed
		return err
	}

	// Return nil if the subscription plan was deleted successfully
	return nil
}
func (u *adminUseCase) GetSubscriptionPlans() ([]domain.SubscriptionPlan, error) {
	// Call the repository function to fetch the list of subscription plans
	plans, err := u.adminRepository.GetSubscriptionPlans()
	if err != nil {
		// Handle any error from the repository, you can log or perform additional actions as needed
		return nil, err
	}

	// Return the list of subscription plans
	return plans, nil
}
