package usecase

import (
	"main/pkg/helper"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"mime/multipart"
)

// UseCase is a struct representing the video use case.
type VideoUseCase struct {
	videoRepo interfaces.VideoRepository
}

// NewVideoUseCase creates a new instance of the video use case.
func NewVideoUseCase(videoRepo interfaces.VideoRepository) services.VideoUseCase {
	return &VideoUseCase{
		videoRepo: videoRepo,
	}
}

// UploadVideo uploads a video file, encodes it, adds it to S3, and stores details in the database.
func (uc *VideoUseCase) UploadVideo(userID int, categoryID int, title, description string, file *multipart.FileHeader) error {
	// Encode video
	videoData, err := helper.EncodeVideo(file)
	if err != nil {
		return err
	}

	// Upload video to S3
	videoURL, err := helper.AddVideoToS3(videoData)
	if err != nil {
		return err
	}

	// Store video details in the database
	if err := uc.videoRepo.UploadVideo(userID, categoryID, title, description, videoURL); err != nil {
		return err
	}

	return nil
}

func (uc *VideoUseCase) ListVideos(userID int, page, limit int) ([]models.Video, error) {
	// Call the repository to get the paginated list of videos
	videos, err := uc.videoRepo.ListVideos(userID, page, limit)
	if err != nil {
		return nil, err
	}

	// If needed, you can perform additional business logic or filtering on the videos here

	return videos, nil
}

// EditVideoDetails edits the title and description of a video.
func (uc *VideoUseCase) EditVideoDetails(videoID int, title, description string) error {
	// Validate the input, if necessary

	// Call the repository to edit video details in the database
	if err := uc.videoRepo.EditVideoDetails(videoID, title, description); err != nil {
		return err
	}

	return nil
}

// DeleteVideo deletes video details and URL based on its ID.
func (uc *VideoUseCase) DeleteVideo(videoID int) error {
	// Delete video details from the database
	if err := uc.videoRepo.DeleteVideo(videoID); err != nil {
		return err
	}

	return nil
}
