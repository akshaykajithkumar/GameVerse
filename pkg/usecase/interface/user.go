package interfaces

import "main/pkg/utils/models"

type UserUseCase interface {
	SignUp(user models.UserDetails) (models.TokenUser, error)
	ChangePassword(id int, old string, password string, repassword string) error
	Login(user models.UserLogin) (models.TokenUser, error)
	EditProfile(id int, name, email, username, phone, bio string) error
	GetProfile(id int) (*models.UserProfileResponse, error)
	ReportUser(reporterID, targetID int, reason string) error
}
