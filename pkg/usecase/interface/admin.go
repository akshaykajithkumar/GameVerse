package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (models.TokenAdmin, error)
	GetUsers(page int, limit int) ([]models.UserDetailsAtAdmin, error)
	ToggleBlockUser(id string) error
	GetReports(page, limit int) ([]domain.Reports, error)
}
