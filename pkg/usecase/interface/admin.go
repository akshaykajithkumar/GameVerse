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
	AddSubscriptionPlan(name string, duration int, price float64) error
	DeleteSubscriptionPlan(planID int) error
	GetSubscriptionPlans() ([]domain.SubscriptionPlan, error)
}
