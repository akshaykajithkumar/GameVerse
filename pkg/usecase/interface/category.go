package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
)

type CategoryUseCase interface {
	AddCategory(category string) (domain.Category, error)
	UpdateCategory(current string, new string) (domain.Category, error)
	DeleteCategory(categoryID string) error
	GetCategories(page, limit int) ([]domain.Category, error)
	ListVideosByCategory(categoryID, page, limit int) ([]models.VideoResponses, error)
}
