package interfaces

import (
	"main/pkg/utils/models"
)

type VideoRepository interface {
	UploadVideo(userID int, categoryID int, title, description, url string) error
	ListVideos(userID, page, limit int) ([]models.Video, error)
	EditVideoDetails(videoID int, title, description string) error
	DeleteVideo(videoID int) error
}
