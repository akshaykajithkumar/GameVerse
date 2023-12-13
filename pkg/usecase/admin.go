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

func (ad *adminUseCase) BlockUser(id string) error {

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if !user.Permission {
		return errors.New("already blocked")
	} else {
		user.Permission = false
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

// unblock user
func (ad *adminUseCase) UnBlockUser(id string) error {

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if !user.Permission {
		user.Permission = true
	} else {
		return errors.New("already unblocked")
	}

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
