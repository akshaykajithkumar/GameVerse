package usecase

import (
	"errors"
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
)

type categoryUseCase struct {
	repository interfaces.CategoryRepository
}

func NewCategoryUseCase(repo interfaces.CategoryRepository) services.CategoryUseCase {
	return &categoryUseCase{
		repository: repo,
	}
}

func (Cat *categoryUseCase) AddCategory(category string) (domain.Category, error) {

	productResponse, err := Cat.repository.AddCategory(category)

	if err != nil {
		return domain.Category{}, err
	}

	return productResponse, nil

}

func (Cat *categoryUseCase) UpdateCategory(current string, new string) (domain.Category, error) {

	result, err := Cat.repository.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}

	if !result {
		return domain.Category{}, errors.New("there is no category as you mentioned")
	}

	newcat, err := Cat.repository.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}

	return newcat, err
}

func (Cat *categoryUseCase) DeleteCategory(categoryID string) error {

	err := Cat.repository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil

}

func (Cat *categoryUseCase) GetCategories(page, limit int) ([]domain.Category, error) {
	categories, err := Cat.repository.GetCategories(page, limit)
	if err != nil {
		return []domain.Category{}, err
	}
	return categories, nil
}

// ListVideosByCategory retrieves a list of videos in a specific category based on category ID, page, and limit.
func (uc *categoryUseCase) ListVideosByCategory(categoryID, page, limit int) ([]models.VideoResponses, error) {
	// Validate categoryID
	if categoryID <= 0 {
		return []models.VideoResponses{}, errors.New("invalid category ID")
	}

	// Validate page and limit
	if page < 1 || limit < 1 {
		return []models.VideoResponses{}, errors.New("page and limit must be positive integers")
	}

	// Fetch videos from the repository
	videos, err := uc.repository.ListVideosByCategory(categoryID, page, limit)
	if err != nil {
		return []models.VideoResponses{}, err
	}

	return videos, nil
}
