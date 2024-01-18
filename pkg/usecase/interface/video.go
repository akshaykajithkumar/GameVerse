package interfaces

import (
	"main/pkg/domain"
	"main/pkg/utils/models"
	"mime/multipart"
)

type VideoUseCase interface {
	UploadVideo(userID int, categoryID int, title, description string, file *multipart.FileHeader, tags []string) error
	ListVideos(userID int, page, limit int) ([]models.Video, error)
	EditVideoDetails(videoID int, title, description string) error
	DeleteVideo(videoID int) error
	WatchVideo(userID int, videoId int) (string, error)
	ToggleLikeVideo(userID uint, videoID uint) error
	CommentVideo(userID uint, videoID uint, content string) error
	GetComments(videoID uint) ([]domain.Comment, error)
	AddVideoTags(tags []string) error
	DeleteVideoTagByID(tagID uint) error
	GetVideoTags() ([]domain.Tag, error)
	StoreUserTags(userID int, tagIDs []uint) error
	// RecommendationList(userID int) ([]models.RecommendationListResponse, error)
	RecommendationList(userID int, page, limit int) ([]models.RecommendationListResponse, error)
}
