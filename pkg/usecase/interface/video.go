package interfaces

import (
	"main/pkg/utils/models"
	"mime/multipart"
)

type VideoUseCase interface {
	UploadVideo(userID int, categoryID int, title, description string, file *multipart.FileHeader) error
	ListVideos(userID int, page, limit int) ([]models.Video, error)
	EditVideoDetails(videoID int, title, description string) error
	DeleteVideo(videoID int) error
}
