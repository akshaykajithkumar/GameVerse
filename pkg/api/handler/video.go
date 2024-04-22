package handler

import (
	"main/pkg/helper"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	VideoUseCase services.VideoUseCase
}

func NewVideoHandler(usecase services.VideoUseCase) *VideoHandler {
	return &VideoHandler{
		VideoUseCase: usecase,
	}
}

// // UploadVideo is a handler for uploading video files.
// // @Summary      Upload Video
// // @Description  Upload a video file along with title, description, and category ID
// // @Tags         User
// // @Accept       multipart/form-data
// // @Produce      json
// // @Security     Bearer
// // @Param        VideoFile   formData  file    true  "Video File"
// // @Param        CategoryID   formData  int     true  "Category ID"
// // @Param        Title        formData  string  true  "Title"
// // @Param        Description  formData  string  true  "Description"
// // @Param        tags         formData   array   true    "Video Tags"
// // @Param        exclusive    query      bool    false   "Exclusive Video"
// // @Success      200  {object} response.Response{}
// // @Failure      400  {object} response.Response{}
// // @Router       /users/upload/video [post]
// func (u *VideoHandler) UploadVideo(c *gin.Context) {
// 	userID, err := helper.GetUserID(c)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Parse form data
// 	categoryID, _ := strconv.Atoi(c.PostForm("CategoryID"))
// 	title := c.PostForm("Title")
// 	description := c.PostForm("Description")

// 	// Retrieve the file from the form data
// 	file, err := c.FormFile("VideoFile")
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Error retrieving video file from form", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Retrieve tags from the form data
// 	tags := c.PostFormArray("tags")

// 	// Parse exclusive parameter from query params
// 	exclusive, _ := strconv.ParseBool(c.Query("exclusive"))

// 	// Call the use case to upload the video, passing the exclusive parameter
// 	if err := u.VideoUseCase.UploadVideo(int(userID), categoryID, title, description, file, tags, exclusive); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not upload video", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Customize the response based on your needs
// 	successRes := response.ClientResponse(http.StatusOK, "Video uploaded successfully", nil, nil)
// 	c.JSON(http.StatusOK, successRes)
// }

// ListVideos is a handler for listing videos for a particular user with pagination.
// @Summary      List Videos
// @Description  List videos for a particular user with pagination
// @Tags         User
// @Security     Bearer
// @Param        limit   query   string  true    "Limit per page"
// @Param        page    query   string  true    "Page number"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/videos [get]
func (u *VideoHandler) ListVideos(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Parse limit and page parameters from the query
	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	// Default values if not provided
	limit := 10
	page := 1

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusBadRequest, "Limit parameter not in the right format", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
	}

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusBadRequest, "Page parameter not in the right format", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
	}

	// Call the use case to list videos for the user with pagination
	videos, err := u.VideoUseCase.ListVideos(int(userID), page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not list videos", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Successfully listed videos", videos, nil)
	c.JSON(http.StatusOK, successRes)
}

// EditVideoDetails is a handler for patching video details such as title and description.
// @Summary      Edit Video Details
// @Description  Patch the title and description of a video
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        VideoID      query   int     true  "Video ID"
// @Param        videoDetailsRequest	body	models.EditVideoDetails	true	"Video Details Request"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router        /users/profile/videos/editVideo [patch]
func (u *VideoHandler) EditVideoDetails(c *gin.Context) {
	// Extract video ID from the query parameter
	videoID, err := strconv.Atoi(c.Query("VideoID"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid video ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Parse the request body into an EditVideoDetails struct
	var editVideoDetails models.EditVideoDetails
	if err := c.ShouldBindJSON(&editVideoDetails); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to edit video details
	if err := u.VideoUseCase.EditVideoDetails(videoID, editVideoDetails.Title, editVideoDetails.Description); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not edit video details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully edited video details", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// DeleteVideo is a handler for deleting a video file.
// @Summary      Delete Video
// @Description  Delete a video file based on its ID
// @Tags         User
// @Produce      json
// @Security     Bearer
// @Param        VideoID      query   int     true  "Video ID to be deleted"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/videos/delete [delete]
func (u *VideoHandler) DeleteVideo(c *gin.Context) {
	// Extract video ID from the query parameter
	videoID, err := strconv.Atoi(c.Query("VideoID"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid video ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to delete the video
	if err := u.VideoUseCase.DeleteVideo(videoID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not delete video", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted video", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// // WatchVideo is a handler for watching a specific video.
// // @Summary      Watch Video
// // @Description  Watch a specific video for a particular user
// // @Tags         User
// // @Security     Bearer
// // @Param        creatorID   query   int     true    "Creator ID"
// // @Param        videoID  query   int     true    "Video ID"
// // @Success      200  {object} response.Response{}
// // @Failure      400  {object} response.Response{}
// // @Router       /users/profile/videos/watch [get]
// func (u *VideoHandler) WatchVideo(c *gin.Context) {
// 	userID, err := helper.GetUserID(c)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	creatorIDStr := c.Query("creatorID")
// 	videoIDStr := c.Query("videoID")

// 	// Validate and parse userID parameter
// 	creatorID, err := strconv.Atoi(creatorIDStr)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "userID parameter not in the right format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Validate and parse videoID parameter
// 	videoID, err := strconv.Atoi(videoIDStr)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "videoID parameter not in the right format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Call the use case to watch the video for the user
// 	videoURL, err := u.VideoUseCase.WatchVideo(userID, videoID, creatorID)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not watch video", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Customize the response based on your needs
// 	successRes := response.ClientResponse(http.StatusOK, "follow the link to watch the video", videoURL, nil)
// 	c.JSON(http.StatusOK, successRes)
// }

// ToggleLikeVideo is a handler for toggling the like status of a video.
// @Summary      Toggle Like Video
// @Description  Toggle the like status of a video for the authenticated user
// @Tags         User
// @Security     Bearer
// @Param        videoID   query   string  true    "Video ID"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/videos/like [post]
func (u *VideoHandler) ToggleLikeVideo(c *gin.Context) {
	// Get user ID using helper function
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get video ID from query parameters
	videoIDStr := c.Query("videoID")
	if videoIDStr == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Validate and parse videoID parameter
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to toggle the like status of the video
	if err := u.VideoUseCase.ToggleLikeVideo(uint(userID), uint(videoID)); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not toggle like status", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Like status toggled successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// CommentVideoHandler is a handler for commenting on a video.
// @Summary      Comment on Video
// @Description  Add a new comment to a video
// @Tags         User
// @Security     Bearer
// @Param        videoID   query   string  true    "Video ID"
// @Param        content   query   string  true    "Comment content"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/videos/comment [post]
func (u *VideoHandler) CommentVideoHandler(c *gin.Context) {
	// Get user ID using helper function
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get video ID from query parameters
	videoIDStr := c.Query("videoID")
	if videoIDStr == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Validate and parse videoID parameter
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get content from query parameters
	content := c.Query("content")
	if content == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Content parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to add a comment to the video
	if err := u.VideoUseCase.CommentVideo(uint(userID), uint(videoID), content); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add comment", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Comment added successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetCommentsHandler is a handler for retrieving comments for a video.
// @Summary      Get Comments
// @Description  Retrieve comments for a specific video
// @Tags         User
// @Security     Bearer
// @Param        videoID   query   uint  true    "Video ID"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/videos/comments [get]
func (u *VideoHandler) GetCommentsHandler(c *gin.Context) {
	// Get video ID from query parameters
	videoIDStr := c.Query("videoID")
	if videoIDStr == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Validate and parse videoID parameter
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "VideoID parameter not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to get comments for the video
	comments, err := u.VideoUseCase.GetComments(uint(videoID))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not retrieve comments", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Comments retrieved successfully", comments, nil)
	c.JSON(http.StatusOK, successRes)
}

// AddTagsHandler is a handler for adding tags to the database.
// @Summary      Add Tags
// @Description  Add tags to the database
// @Tags         Admin
// @Security     Bearer
// @Param        tags   query   string  true    "Comma-separated list of tags"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /admin/addtags [post]
func (u *VideoHandler) AddTagsHandler(c *gin.Context) {
	// Get tags from query parameters
	tagsParam := c.Query("tags")
	if tagsParam == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Tags parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Split the comma-separated tags into a slice
	tags := strings.Split(tagsParam, ",")

	// Call the use case to add tags to the database
	if err := u.VideoUseCase.AddVideoTags(tags); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add tags", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Tags added successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// DeleteTagHandler is a handler for deleting a tag from the database based on a tag ID.
// @Summary      Delete Tag
// @Description  Delete a tag from the database based on a tag ID
// @Tags         Admin
// @Security     Bearer
// @Param        tagID   query   uint  true    "Tag ID to delete"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /admin/deletetags [delete]
func (u *VideoHandler) DeleteTagHandler(c *gin.Context) {
	// Get tag ID from query parameters
	tagIDParam := c.Query("tagID")
	if tagIDParam == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "TagID parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Convert tag ID to uint
	tagID, err := strconv.ParseUint(tagIDParam, 10, 64)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "TagID parameter not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to delete the tag from the database based on the tag ID
	if err := u.VideoUseCase.DeleteVideoTagByID(uint(tagID)); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not delete tag", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Tag deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetTagsHandler is a handler for getting a list of tags.
// @Summary      Get Tags
// @Description  Get a list of tags
// @Tags         Admin
// @Security     Bearer
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /admin/tags [get]
func (u *VideoHandler) GetTagsHandler(c *gin.Context) {
	// Call the use case to get a list of tags
	tags, err := u.VideoUseCase.GetVideoTags()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get tags", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Tags retrieved successfully", tags, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetTagsHandler is a handler for getting a list of tags.
// @Summary      Get Tags
// @Description  Get a list of tags
// @Tags         User
// @Security     Bearer
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/tags [get]
func (u *VideoHandler) GetTagsForUserHandler(c *gin.Context) {
	// Call the use case to get a list of tags
	tags, err := u.VideoUseCase.GetVideoTags()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get tags", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Tags retrieved successfully", tags, nil)
	c.JSON(http.StatusOK, successRes)
}

// StoreUserTags is a handler for storing multiple tags for a specific user.
// @Summary      Store User Tags
// @Description  Store multiple tags for a specific user
// @Tags         User
// @Security     Bearer
// @Param        tagIDs   query  string true   "Comma-separated list of tag IDs (e.g., {1,2})"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/selectTags [post]
func (u *VideoHandler) StoreUserTags(c *gin.Context) {
	// Get userID using helper function
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get tagIDs from query parameter
	tagIDsParam := c.Query("tagIDs")
	if tagIDsParam == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "TagIDs parameter is required", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Remove curly braces from the tagIDsParam
	tagIDsParam = strings.Trim(tagIDsParam, "{}")

	// Split the comma-separated tagIDs into a slice of strings
	tagIDsStr := strings.Split(tagIDsParam, ",")

	// Convert tagIDs from strings to uint
	var tagIDs []uint
	for _, tagIDStr := range tagIDsStr {
		tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid tagID", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
		tagIDs = append(tagIDs, uint(tagID))
	}

	// Call the use case to store user tags
	if err := u.VideoUseCase.StoreUserTags(userID, tagIDs); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not store user tags", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "User tags stored successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// // @Summary Recommendation List
// // @Description Generate a recommendation list of videos for the authenticated user based on tags.
// // @Tags User
// // @Security Bearer
// // @Success 200 {array} models.RecommendationListResponse
// // @Failure 400 {object} response.Response
// // @Router /users/profile/videos/recommendation [get]
// func (u *VideoHandler) RecommendationList(c *gin.Context) {
// 	userID, err := helper.GetUserID(c)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Call the use case to get the recommendation list for the user
// 	recommendations, err := u.VideoUseCase.RecommendationList(int(userID))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get recommendation list", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	// Customize the response based on your needs
// 	successRes := response.ClientResponse(http.StatusOK, "Recommendation list retrieved successfully", recommendations, nil)
// 	c.JSON(http.StatusOK, successRes)
// }

// @Summary Recommendation List
// @Description Generate a recommendation list of videos for the authenticated user based on tags.
// @Tags User
// @Security Bearer
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} models.RecommendationListResponse
// @Failure 400 {object} response.Response
// @Router /users/profile/videos/recommendation [get]
func (u *VideoHandler) RecommendationList(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Parse query parameters for pagination
	page, limit := parsePaginationParams(c)

	// Call the use case to get the paginated recommendation list for the user
	recommendations, err := u.VideoUseCase.RecommendationList(int(userID), page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get recommendation list", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Recommendation list retrieved successfully", recommendations, nil)
	c.JSON(http.StatusOK, successRes)
}

func parsePaginationParams(c *gin.Context) (int, int) {
	// Default values for page and limit
	page := 1
	limit := 10

	// Parse query parameters for page and limit
	if pageStr, ok := c.GetQuery("page"); ok {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if limitStr, ok := c.GetQuery("limit"); ok {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	return page, limit
}

// ListVideos is a handler for searching and listing videos with sorting and pagination.
// @Summary      List/Search Videos
// @Description  List/Search videos with sorting and pagination
// @Tags         User
// @Security     Bearer
// @Param        limit   query   int     false   "Limit per page"
// @Param        page    query   int     false   "Page number"
// @Param        sort    query   string  false   "Sort order (upload_time, views, likes)"
// @Param        order   query   string  false   "Order (asc, desc)"
// @Param        search  query   string  false   "Search term"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/videos [get]
func (u *VideoHandler) ListtVideos(c *gin.Context) {
	// Parse query parameters
	limit, page, sort, order, search := parseListVideosQueryParams(c)

	// Call the use case to search and list videos with sorting and pagination
	videos, err := u.VideoUseCase.ListtVideos(page, limit, sort, order, search)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not list videos", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Customize the response based on your needs
	successRes := response.ClientResponse(http.StatusOK, "Successfully listed videos", videos, nil)
	c.JSON(http.StatusOK, successRes)
}

// parseListVideosQueryParams parses and validates query parameters for ListVideos
func parseListVideosQueryParams(c *gin.Context) (int, int, string, string, string) {
	// Default values
	limit := 10
	page := 1
	sort := ""
	order := "asc"
	search := ""

	// Parse limit and page parameters from the query
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
	sort = c.Query("sort")
	order = c.Query("order")
	search = c.Query("search")

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	return limit, page, sort, order, search
}
