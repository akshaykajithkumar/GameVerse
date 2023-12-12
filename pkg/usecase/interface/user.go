package interfaces

import "main/pkg/utils/models"

type UserUseCase interface {
	SignUp(user models.UserDetails) (models.TokenUser, error)
}
