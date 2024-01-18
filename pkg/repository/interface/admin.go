package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	GetUsers(page int, limit int) ([]models.UserDetailsAtAdmin, error)
	GetUserByID(id string) (domain.User, error)
	UpdateBlockUserByID(user domain.User) error
	GetReports(page, limit int) ([]domain.Reports, error)
	AddSubscriptionPlan(name string, duration int, price float64) error
	DeleteSubscriptionPlan(planID int) error
	GetSubscriptionPlans() ([]domain.SubscriptionPlan, error)
}
