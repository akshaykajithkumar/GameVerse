package repository

import (
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"

	"gorm.io/gorm"
)

// VideoRepository is a struct representing the video repository.
type VideoRepository struct {
	DB *gorm.DB
}

// NewVideoRepository creates a new instance of the video repository.
func NewVideoRepository(db *gorm.DB) interfaces.VideoRepository {
	return &VideoRepository{
		DB: db,
	}
}

// UploadVideo uploads video details to the database.
func (vr *VideoRepository) UploadVideo(userID int, categoryID int, title, description, url string) error {
	video := domain.Video{
		UserID:      uint(userID),
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
		URL:         url,
	}

	if err := vr.DB.Create(&video).Error; err != nil {
		return err
	}

	return nil
}

func (vr *VideoRepository) ListVideos(userID, page, limit int) ([]models.Video, error) {
	var videos []models.Video

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	// Query the database with pagination
	if err := vr.DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

// EditVideoDetails updates the title and description of a video in the database.
func (vr *VideoRepository) EditVideoDetails(videoID int, title, description string) error {
	// Find the video by ID
	var video domain.Video
	if err := vr.DB.First(&video, videoID).Error; err != nil {
		return err
	}

	// Update the video details
	video.Title = title
	video.Description = description

	// Save the changes back to the database
	if err := vr.DB.Save(&video).Error; err != nil {
		return err
	}

	return nil
}

// DeleteVideo deletes video details from the database based on video ID.
func (vr *VideoRepository) DeleteVideo(videoID int) error {
	if err := vr.DB.Where("id = ?", videoID).Delete(&domain.Video{}).Error; err != nil {
		return err
	}

	return nil
}
