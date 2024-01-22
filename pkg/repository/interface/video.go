package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
)

type VideoRepository interface {
	UploadVideo(userID int, categoryID int, title, description, url string, tags []string, exclusive bool) (uint, error)
	ListVideos(userID, page, limit int) ([]models.Video, error)
	EditVideoDetails(videoID int, title, description string) error
	DeleteVideo(videoID int) error
	IncrementVideoViews(videoID int) error
	IsLikedByUser(userID uint, videoID uint) bool
	UnlikeVideo(userID uint, videoID uint) error
	LikeVideo(userID uint, videoID uint) error
	CreateComment(comment *domain.Comment) error
	GetCommentsByVideoID(videoID uint) ([]domain.Comment, error)
	AddTags(tags []string) error
	DeleteTagByID(tagID uint) error
	GetTags() ([]domain.Tag, error)
	StoreUserTags(userTags []domain.UserTags) error
	GetAllVideos() ([]domain.Video, error)
	// GetUserTags(userID int) ([]domain.UserTags, error)
	GetUserTags(userID int) ([]string, error)
	UpdateVideoLikesCount(videoID uint) error
	GetVideoTagsByVideoID(videoID uint) ([]string, error)
	IsUserSubscribed(userID int, creatorID int) (bool, error)
	IsVideoExclusive(videoID int) (bool, error)
}
