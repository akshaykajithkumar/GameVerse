package repository

import (
	"errors"
	"main/pkg/domain"
	interfaces "main/pkg/repository/interface"
	"main/pkg/utils/models"
	"time"

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
// func (vr *VideoRepository) UploadVideo(userID int, categoryID int, title, description, url string) error {
// 	video := domain.Video{
// 		UserID:      uint(userID),
// 		CategoryID:  categoryID,
// 		Title:       title,
// 		Description: description,
// 		URL:         url,
// 	}

// 	if err := vr.DB.Create(&video).Error; err != nil {
// 		return err
// 	}

//		return nil
//	}
//
// UploadVideo stores video details in the database, including tags.

// func (vr *VideoRepository) UploadVideo(userID int, categoryID int, title, description, url string, tags []string) (uint, error) {
// 	// Create a new Video instance
// 	video := domain.Video{
// 		UserID:      uint(userID),
// 		CategoryID:  categoryID,
// 		Title:       title,
// 		Description: description,
// 		URL:         url,
// 		// Tags:        tags,
// 	}

// 	if err := vr.DB.Create(&video).Error; err != nil {
// 		return 0, err
// 	}

//		// Assuming that video.ID is the auto-generated ID of the newly created video
//		return video.ID, nil
//	}
func (vr *VideoRepository) UploadVideo(userID int, categoryID int, title, description, url string, tags []string, exclusive bool) (uint, error) {
	// Create a new Video instance
	video := domain.Video{
		UserID:      uint(userID),
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
		URL:         url,
		Exclusive:   exclusive,
		CreatedAt:   time.Now(), // Set the creation time to the current time
	}

	if err := vr.DB.Create(&video).Error; err != nil {
		return 0, err
	}

	// Insert tags into VideoTags table for the newly created video
	for _, tag := range tags {
		videoTag := domain.VideoTags{
			UserID:  userID,
			VideoID: video.ID,
			Tag:     tag,
		}

		if err := vr.DB.Create(&videoTag).Error; err != nil {
			return 0, err
		}
	}

	// Return the auto-generated ID of the newly created video
	return video.ID, nil
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

// IncrementVideoViews increments the views count for a specific video.
func (repo *VideoRepository) IncrementVideoViews(videoID int) error {
	var video domain.Video
	err := repo.DB.Model(&video).Where("id = ?", videoID).Update("views", gorm.Expr("views + ?", 1)).Error
	if err != nil {
		// Handle the error
		return err
	}
	return nil
}
func (vr *VideoRepository) IsLikedByUser(userID uint, videoID uint) bool {
	var likeCount int64
	err := vr.DB.Model(&domain.VideoLikes{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&likeCount).Error
	if err != nil {
		// Handle the error (logging, return false, etc.)
		return false
	}

	// If likeCount is greater than 0, the user has already liked the video
	return likeCount > 0
}

func (vr *VideoRepository) UnlikeVideo(userID uint, videoID uint) error {
	return vr.DB.Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&domain.VideoLikes{}).
		Error
}

func (vr *VideoRepository) LikeVideo(userID uint, videoID uint) error {
	like := &domain.VideoLikes{
		UserID:  userID,
		VideoID: videoID,
	}
	return vr.DB.Create(like).Error
}
func (vr *VideoRepository) UpdateVideoLikesCount(videoID uint) error {
	var likeCount int64
	if err := vr.DB.Model(&domain.VideoLikes{}).Where("video_id = ?", videoID).Count(&likeCount).Error; err != nil {
		return err
	}

	// Update the Video table with the new like count
	if err := vr.DB.Model(&domain.Video{}).Where("id = ?", videoID).Update("likes", likeCount).Error; err != nil {
		return err
	}

	return nil
}

// CreateComment adds a new comment to the repository.
func (vr *VideoRepository) CreateComment(comment *domain.Comment) error {
	return vr.DB.Create(comment).Error
}

// GetCommentsByVideoID retrieves comments for a specific video from the database.
func (vr *VideoRepository) GetCommentsByVideoID(videoID uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	// Assuming you have a 'comments' table in your database with a 'video_id' column
	if err := vr.DB.Where("video_id = ?", videoID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

// AddTags adds multiple tags to the database.
func (vr *VideoRepository) AddTags(tags []string) error {
	// Create a slice to store individual tags
	var tagRecords []domain.Tag

	// Iterate through the tags and create Tag instances
	for _, tag := range tags {
		tagRecord := domain.Tag{
			Tag: tag,
		}
		tagRecords = append(tagRecords, tagRecord)
	}

	// Insert all tags in a single database transaction
	if err := vr.DB.Create(&tagRecords).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTagByID deletes a tag from the database based on a tag ID.
func (vr *VideoRepository) DeleteTagByID(tagID uint) error {
	// Delete the tag based on the tag ID
	if err := vr.DB.Where("id = ?", tagID).Delete(&domain.Tag{}).Error; err != nil {
		return err
	}

	return nil
}

// GetTags retrieves a list of tags from the database.
func (vr *VideoRepository) GetTags() ([]domain.Tag, error) {
	var tags []domain.Tag

	// Retrieve all tags from the database
	if err := vr.DB.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// StoreUserTags stores multiple user tags in the UserTags table.
func (vr *VideoRepository) StoreUserTags(userTags []domain.UserTags) error {
	// Insert the user tags into the UserTags table
	if err := vr.DB.Create(&userTags).Error; err != nil {
		return err
	}

	return nil
}

// GetAllVideos retrieves all videos.
//
//	func (vr *VideoRepository) GetAllVideos() ([]domain.Video, error) {
//		var videos []domain.Video
//		if err := vr.DB.Find(&videos).Error; err != nil {
//			return nil, err
//		}
//		return videos, nil
//	}
func (vr *VideoRepository) GetAllVideos() ([]domain.Video, error) {
	var videos []domain.Video

	// Use raw SQL query to retrieve only necessary fields from videos
	if err := vr.DB.Raw("SELECT id, user_id, title, description, url, category_id, likes, views FROM videos").Scan(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

// func (vr *VideoRepository) GetUserTags(userID int) ([]domain.UserTags, error) {
// 	var userTags []domain.UserTags

// 	// Query the database to retrieve all UserTags associated with the user
// 	if err := vr.DB.Where("user_id = ?", userID).Find(&userTags).Error; err != nil {
// 		return nil, err
// 	}

//		return userTags, nil
//	}
func (vr *VideoRepository) GetUserTags(userID int) ([]string, error) {
	var tags []string

	// Define the raw SQL query to select distinct tags where UserID is equal to the specified parameter
	sqlQuery := "SELECT DISTINCT tags.tag FROM user_tags INNER JOIN tags ON user_tags.tag_id = tags.id WHERE user_tags.user_id = ?"

	// Execute the raw SQL query
	if err := vr.DB.Raw(sqlQuery, userID).Pluck("tag", &tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}
func (vr *VideoRepository) GetVideoTagsByVideoID(videoID uint) ([]string, error) {
	var tags []string

	// Define the raw SQL query to select tags where VideoID is equal to the specified parameter
	sqlQuery := "SELECT tag FROM video_tags WHERE video_id = ?"

	// Execute the raw SQL query
	if err := vr.DB.Raw(sqlQuery, videoID).Pluck("tag", &tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// IsVideoExclusive checks if a video is marked as exclusive.
func (vr *VideoRepository) IsVideoExclusive(videoID int) (bool, error) {
	var video domain.Video
	err := vr.DB.Where("id = ?", videoID).First(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Video not found
			return false, errors.New("video not found")
		}
		// Error occurred while querying the database
		return false, err
	}

	return video.Exclusive, nil
}

// IsUserSubscribed checks if a user is subscribed to a specific creator.
// func (vr *VideoRepository) IsUserSubscribed(userID int, creatorID int) (bool, error) {
// 	var subscription domain.SubscriptionList
// 	err := vr.DB.Where("user_id = ? AND creator_id = ? AND is_active = ?", userID, creatorID, true).First(&subscription).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// User is not subscribed
// 			return false, nil
// 		}
// 		// Error occurred while querying the database
// 		return false, err
// 	}

//		// User is subscribed
//		return true, nil
//	}
func (vr *VideoRepository) IsUserSubscribed(userID int, creatorID int) (bool, error) {
	var count int
	err := vr.DB.Raw("SELECT COUNT(*) FROM subscription_lists WHERE user_id = ? AND creator_id = ? AND is_active = true", userID, creatorID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	// If the count is greater than 0, the user is subscribed; otherwise, the user is not subscribed
	return count > 0, nil
}

// ListVideos is a repository method for searching and listing videos with sorting and pagination.
func (vr *VideoRepository) ListtVideos(page, limit int, sort, order, search string) ([]models.Video, error) {
	var videos []models.Video

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	// Query the database with sorting and pagination
	query := vr.DB

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	switch sort {
	case "views":
		query = query.Order("views " + order)
	case "likes":
		query = query.Order("likes " + order)
	case "created_at":
		query = query.Order("created_at " + order)
	default:
		// Default sorting by title in ascending order
		query = query.Order("title ASC")
	}

	if err := query.Offset(offset).Limit(limit).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}
