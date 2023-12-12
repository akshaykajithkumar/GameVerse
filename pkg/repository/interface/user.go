package interfaces

import "main/pkg/utils/models"

type UserRepository interface {
	SignUp(user models.UserDetails) (models.UserResponse, error)
	CheckUserAvailability(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserResponse, error)
}
