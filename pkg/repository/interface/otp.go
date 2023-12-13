package interfaces

import "main/pkg/utils/models"

type OtpRepository interface {
	FindUserByMobileNumber(phone string) bool
	UserDetailsUsingPhone(phone string) (models.UserResponse, error)
	UpdateUserPermissionByPhone(phone string) error
}
