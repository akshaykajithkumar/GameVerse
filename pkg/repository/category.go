package repository

import (
	"errors"
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"
	"strconv"

	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) interfaces.CategoryRepository {
	return &categoryRepository{DB}
}

func (p *categoryRepository) AddCategory(cat string) (domain.Category, error) {

	var b string
	err := p.DB.Raw("INSERT INTO categories (category) VALUES (?) RETURNING category", cat).Scan(&b).Error
	if err != nil {
		return domain.Category{}, err
	}

	var categoryResponse domain.Category
	err = p.DB.Raw(`
	SELECT
		p.id,
		p.category
		FROM
			categories p
		WHERE
			p.category = ?
			`, b).Scan(&categoryResponse).Error

	if err != nil {
		return domain.Category{}, err
	}

	return categoryResponse, nil

}

func (p *categoryRepository) CheckCategory(current string) (bool, error) {
	var i int
	err := p.DB.Raw("SELECT COUNT(*) FROM categories WHERE category=?", current).Scan(&i).Error
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, err
	}

	return true, err
}

func (p *categoryRepository) UpdateCategory(current, new string) (domain.Category, error) {

	// Check the database connection
	if p.DB == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}

	// Update the category
	if err := p.DB.Exec("UPDATE categories SET category = ? WHERE category = ?", new, current).Error; err != nil {
		return domain.Category{}, err
	}

	// Retrieve the updated category
	var newcat domain.Category
	if err := p.DB.First(&newcat, "category = ?", new).Error; err != nil {
		return domain.Category{}, err
	}

	return newcat, nil
}

func (c *categoryRepository) DeleteCategory(categoryID string) error {
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return errors.New("converting into integer not happened")
	}

	result := c.DB.Exec("DELETE FROM categories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (c *categoryRepository) GetCategories(page, limit int) ([]domain.Category, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	var categories []domain.Category

	if err := c.DB.Raw("select id,category from categories limit ? offset ?", limit, offset).Scan(&categories).Error; err != nil {
		return []domain.Category{}, err
	}

	return categories, nil
}

// func (vr *categoryRepository) ListVideosByCategory(categoryID, page, limit int) (models.VideoResponses, error) {
// 	// Validate categoryID
// 	if categoryID <= 0 {
// 		return models.VideoResponses{}, errors.New("invalid category ID")
// 	}

// 	// Validate page and limit
// 	if page < 1 || limit < 1 {
// 		return models.VideoResponses{}, errors.New("page and limit must be positive integers")
// 	}

// 	// Calculate offset for pagination
// 	offset := (page - 1) * limit

// 	var videos models.VideoResponses

// 	// Fetch videos from the database based on category ID, page, and limit
// 	if err := vr.DB.Table("videos").
// 		Select("id, user_id, title, description, url, category_id").
// 		Where("category_id = ?", categoryID).
// 		Offset(offset).
// 		Limit(limit).
// 		Scan(&videos).Error; err != nil {
// 		return models.VideoResponses{}, err
// 	}

//		return videos, nil
//	}

func (vr *categoryRepository) ListVideosByCategory(categoryID, page, limit int) ([]models.VideoResponses, error) {
	// Validate categoryID
	if categoryID <= 0 {
		return nil, errors.New("invalid category ID")
	}

	// Validate page and limit
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit must be positive integers")
	}

	// Calculate offset for pagination
	offset := (page - 1) * limit

	var videos []models.VideoResponses

	// Fetch videos from the database based on category ID, page, and limit
	if err := vr.DB.Raw("SELECT id, user_id, title, description, url, category_id FROM videos WHERE category_id = ? OFFSET ? LIMIT ?", categoryID, offset, limit).
		Scan(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}
