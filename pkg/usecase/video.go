package usecase

import (
	// "fmt"
	"sort"

	// "github.com/agnivade/levenshtein"

	"context"
	"errors"
	"log"
	"main/pkg/domain"
	"main/pkg/helper"
	"main/pkg/helper/kafka"
	interfaces "main/pkg/repository/interface"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"mime/multipart"
	"time"

	"github.com/agnivade/levenshtein"
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
func (uc *VideoUseCase) UploadVideo(userID int, categoryID int, title, description string, file *multipart.FileHeader, tags []string, exclusive bool) error {
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
	videoID, err := uc.videoRepo.UploadVideo(userID, categoryID, title, description, videoURL, tags, exclusive)
	if err != nil {
		return err
	}

	if err := kafka.ProduceToKafka(userID, int(videoID), videoURL); err != nil {
		// Handle the Kafka produce error appropriately
		log.Println("Error producing to Kafka:", err)
		// You might choose to return the error or handle it differently
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

// WatchVideo watches a video for the specified user and video ID.
// func (uc *VideoUseCase) WatchVideo(userID int, videoID int) (string, error) {
// 	kafkaTopic := "video_stream"

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	urlChannel := make(chan string)
// 	defer close(urlChannel)

// 	go func() {
// 		// Start Kafka consumer in a goroutine
// 		kafka.ConsumeURLFromKafkaForUser(ctx, kafkaTopic, userID, videoID, urlChannel)
// 	}()

// 	// Wait for video URL from Kafka or timeout after 30 seconds
// 	select {
// 	case videoURL := <-urlChannel:
// 		// Increment the Views count for the watched video
// 		if err := uc.videoRepo.IncrementVideoViews(videoID); err != nil {
// 			// Handle error if needed
// 			return "", err
// 		}

// 		// Optionally, perform additional actions with the videoURL if needed
// 		return videoURL, nil

//		case <-time.After(30 * time.Second):
//			// Timeout case to avoid blocking indefinitely
//			return "", errors.New("timed out waiting for video URL from Kafka")
//		}
//	}
func (uc *VideoUseCase) WatchVideo(userID, videoID, creatorID int) (string, error) {
	// Check if the video is exclusive
	isExclusive, err := uc.videoRepo.IsVideoExclusive(videoID)
	if err != nil {
		// Handle error if needed
		return "", err
	}

	if isExclusive {
		// Check if the user is subscribed to the creator
		isSubscribed, err := uc.videoRepo.IsUserSubscribed(userID, creatorID)
		if err != nil {
			// Handle error if needed
			return "", err
		}

		if !isSubscribed {
			return "", errors.New("user is not subscribed to the creator")
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	urlChannel := make(chan string)

	go func() {
		// Start Kafka consumer in a goroutine
		kafka.ConsumeURLFromKafkaForUser(ctx, userID, videoID, urlChannel)
	}()

	select {
	case videoURL := <-urlChannel:
		// Increment the Views count for the watched video
		if err := uc.videoRepo.IncrementVideoViews(videoID); err != nil {
			// Handle error if needed
			return "", err
		}

		// Optionally, perform additional actions with the videoURL if needed
		return videoURL, nil

	case <-ctx.Done():
		// Context done, check if it's due to a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.New("timed out waiting for video URL from Kafka")
		}
		// Context done for other reasons (possibly due to an error in consumer)
		return "", ctx.Err()
	}
}
func (uc *VideoUseCase) ToggleLikeVideo(userID uint, videoID uint) error {
	// Check if the user has already liked the video
	likedByUser := uc.videoRepo.IsLikedByUser(userID, videoID)

	// Toggle the like status directly in the repository
	if likedByUser {
		// Unlike the video
		if err := uc.videoRepo.UnlikeVideo(userID, videoID); err != nil {
			return err
		}
	} else {
		// Like the video
		if err := uc.videoRepo.LikeVideo(userID, videoID); err != nil {
			return err
		}
	}
	// Update the Likes count in the Video table
	if err := uc.videoRepo.UpdateVideoLikesCount(videoID); err != nil {
		return err
	}
	return nil
}

// CommentVideo adds a new comment to a video.
func (uc *VideoUseCase) CommentVideo(userID uint, videoID uint, content string) error {
	// Create a new comment
	comment := &domain.Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}

	// Add the comment to the repository
	if err := uc.videoRepo.CreateComment(comment); err != nil {
		return err
	}

	return nil
}

// GetComments retrieves comments for a specific video.
func (uc *VideoUseCase) GetComments(videoID uint) ([]domain.Comment, error) {
	// Call the repository function to get comments for the video
	comments, err := uc.videoRepo.GetCommentsByVideoID(videoID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// AddVideoTags adds multiple tags to a video.
func (uc *VideoUseCase) AddVideoTags(tags []string) error {
	// Call the repository function to add tags to the database
	if err := uc.videoRepo.AddTags(tags); err != nil {
		return err
	}

	return nil
}

// DeleteVideoTagByID deletes a tag from the database based on a tag ID.
func (uc *VideoUseCase) DeleteVideoTagByID(tagID uint) error {
	// Call the repository function to delete the tag from the database based on the tag ID
	if err := uc.videoRepo.DeleteTagByID(tagID); err != nil {
		return err
	}

	return nil
}

// GetVideoTags retrieves a list of tags.
func (uc *VideoUseCase) GetVideoTags() ([]domain.Tag, error) {
	// Call the repository function to get a list of tags
	tags, err := uc.videoRepo.GetTags()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// StoreUserTags stores multiple tags for a specific user in the UserTags table.
func (uc *VideoUseCase) StoreUserTags(userID int, tagIDs []uint) error {
	// Create a slice of UserTags instances for each tag ID
	var userTags []domain.UserTags
	for _, tagID := range tagIDs {
		userTags = append(userTags, domain.UserTags{
			UserID: userID,
			TagID:  tagID,
		})
	}

	// Call the repository function to store the user tags
	if err := uc.videoRepo.StoreUserTags(userTags); err != nil {
		return err
	}

	return nil
}

// RecommendationList generates a recommendation list of video IDs for a user based on fuzzy matching of tags.

// RecommendationList generates a recommendation list of video details for a user based on fuzzy matching of tags.
// func (uc *VideoUseCase) RecommendationList(userID int) ([]models.RecommendationListResponse, error) {
// 	// Retrieve user tags from the UserTags table
// 	userTags, err := uc.videoRepo.GetUserTags(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Retrieve all videos from the Video table
// 	allVideos, err := uc.videoRepo.GetAllVideos()

// 	fmt.Println(allVideos)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a map to store video scores based on fuzzy matching
// 	videoScores := make(map[uint]int)

// 	// Iterate over each video
// 	for _, video := range allVideos {
// 		// Calculate the total score for the video based on fuzzy matching of tags
// 		totalScore := 0

// 		// Iterate over each tag of the video
// 		for _, videoTag := range video.Tags {
// 			// Iterate over each user tag
// 			for _, userTag := range userTags {
// 				// Assuming a threshold of 2 for a match, adjust as needed
// 				if distance := levenshtein.ComputeDistance(videoTag, userTag.Tag.Tag); distance <= 2 {
// 					// Increase the total score if there is a match
// 					totalScore++
// 				}
// 			}
// 		}

// 		// Store the total score for the video
// 		videoScores[video.ID] = totalScore
// 	}

// 	// Sort video IDs based on scores (descending order)
// 	sortedVideoIDs := sortVideoIDsByScore(allVideos, videoScores)

// 	// Prepare the response with video details and scores
// 	var recommendations []models.RecommendationListResponse
// 	for _, videoID := range sortedVideoIDs {
// 		for _, video := range allVideos {
// 			if video.ID == videoID {
// 				recommendations = append(recommendations, models.RecommendationListResponse{
// 					ID:          video.ID,
// 					UserID:      video.UserID,
// 					Title:       video.Title,
// 					Description: video.Description,
// 					URL:         video.URL,
// 				})
// 				break
// 			}
// 		}
// 	}

// 		return recommendations, nil
// 	}

// Assuming you have a distance calculation package (levenshtein) imported
// Now, in your RecommendationList function:
// func (uc *VideoUseCase) RecommendationList(userID int) ([]models.RecommendationListResponse, error) {
// 	// Retrieve user tags from the UserTags table
// 	userTags, err := uc.videoRepo.GetUserTags(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a hashmap to store user tags for quick lookup
// 	userTagsMap := make(map[string]struct{})
// 	for _, userTag := range userTags {
// 		userTagsMap[userTag] = struct{}{}
// 	}

// 	// Retrieve all videos from the Video table
// 	allVideos, err := uc.videoRepo.GetAllVideos()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a map to store video scores based on fuzzy matching
// 	videoScores := make(map[uint]int)

// 	// Iterate over each video
// 	for _, video := range allVideos {
// 		// Calculate the total score for the video based on fuzzy matching of tags
// 		totalScore := 0

// 		// Retrieve tags for the current video
// 		currentVideoTags, err := uc.videoRepo.GetVideoTagsByVideoID(video.ID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Iterate over each tag of the current video
// 		for _, currentVideoTag := range currentVideoTags {
// 			// Check if the tag is in the user's tags using the hashmap
// 			if _, exists := userTagsMap[currentVideoTag]; exists {
// 				// Increase the total score if there is a match
// 				totalScore++
// 			}
// 		}

// 		// Store the total score for the video
// 		videoScores[video.ID] = totalScore
// 	}

// 	// Sort video IDs based on scores (descending order)
// 	sortedVideoIDs := sortVideoIDsByScore(allVideos, videoScores)

// 	// Prepare the response with video details and scores
// 	var recommendations []models.RecommendationListResponse
// 	for _, videoID := range sortedVideoIDs {
// 		for _, video := range allVideos {
// 			if video.ID == videoID {
// 				recommendations = append(recommendations, models.RecommendationListResponse{
// 					ID:          video.ID,
// 					UserID:      video.UserID,
// 					Title:       video.Title,
// 					Description: video.Description,
// 					URL:         video.URL,
// 				})
// 				break
// 			}
// 		}
// 	}

// 	return recommendations, nil
// }

// Assuming you have a distance calculation package (levenshtein) imported

func (uc *VideoUseCase) RecommendationList(userID int, page, limit int) ([]models.RecommendationListResponse, error) {
	// Retrieve user tags from the UserTags table
	userTags, err := uc.videoRepo.GetUserTags(userID)
	if err != nil {
		return nil, err
	}

	// Retrieve all videos from the Video table
	allVideos, err := uc.videoRepo.GetAllVideos()
	if err != nil {
		return nil, err
	}

	// Create a map to store video scores based on fuzzy matching
	videoScores := make(map[uint]int)

	// Iterate over each video
	for _, video := range allVideos {
		// Calculate the total score for the video based on fuzzy matching of tags
		totalScore := 0

		// Retrieve tags for the current video
		currentVideoTags, err := uc.videoRepo.GetVideoTagsByVideoID(video.ID)
		if err != nil {
			return nil, err
		}

		// Iterate over each tag of the current video
		for _, currentVideoTag := range currentVideoTags {
			// Iterate over each user tag
			for _, userTag := range userTags {
				// Assuming a threshold of 2 for a match, adjust as needed
				if distance := levenshtein.ComputeDistance(currentVideoTag, userTag); distance <= 2 {
					// Increase the total score if there is a match
					totalScore++
				}
			}
		}

		// Store the total score for the video
		videoScores[video.ID] = totalScore
	}

	// Sort video IDs based on scores (descending order)
	sortedVideoIDs := sortVideoIDsByScore(allVideos, videoScores)

	// Paginate the sorted video IDs based on the provided page and limit
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > len(sortedVideoIDs) {
		endIndex = len(sortedVideoIDs)
	}

	// Prepare the response with video details and scores for the paginated range
	var recommendations []models.RecommendationListResponse
	for _, videoID := range sortedVideoIDs[startIndex:endIndex] {
		for _, video := range allVideos {
			if video.ID == videoID {
				recommendations = append(recommendations, models.RecommendationListResponse{
					ID:          video.ID,
					UserID:      video.UserID,
					Title:       video.Title,
					Description: video.Description,
					URL:         video.URL,
				})
				break
			}
		}
	}

	return recommendations, nil
}

// func sortVideoIDsByScore(videos []domain.Video, scores map[uint]int) []uint {
// 	// Create a slice to store video IDs
// 	var videoIDs []uint

// 	// Iterate over the videos and add their IDs to the slice
// 	for _, video := range videos {
// 		videoIDs = append(videoIDs, video.ID)
// 	}

// 	// Create a custom sorting function based on video scores
// 	sort.Slice(videoIDs, func(i, j int) bool {
// 		return scores[videoIDs[i]] > scores[videoIDs[j]]
// 	})

//		return videoIDs
//	}
func sortVideoIDsByScore(videos []domain.Video, scores map[uint]int) []uint {
	// Create a slice to store video IDs
	var videoIDs []uint

	// Iterate over the videos and add their IDs to the slice only if the score is greater than zero
	for _, video := range videos {
		if score, exists := scores[video.ID]; exists && score > 0 {
			videoIDs = append(videoIDs, video.ID)
		}
	}

	// Create a custom sorting function based on video scores
	sort.Slice(videoIDs, func(i, j int) bool {
		return scores[videoIDs[i]] > scores[videoIDs[j]]
	})

	return videoIDs
}
