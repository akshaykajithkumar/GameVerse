package handler

import (
	"main/pkg/helper"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"
	"strconv"

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

// UploadVideo is a handler for uploading video files.
// @Summary      Upload Video
// @Description  Upload a video file along with title, description, and category ID
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param        VideoFile   formData  file    true  "Video File"
// @Param        CategoryID   formData  int     true  "Category ID"
// @Param        Title        formData  string  true  "Title"
// @Param        Description  formData  string  true  "Description"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/upload/video [post] // Adjust the route to match the one in UserRoutes
func (u *VideoHandler) UploadVideo(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Retrieve file from form
	file, err := c.FormFile("VideoFile")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Error retrieving video file from form", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Extract additional parameters from the form
	categoryID, _ := strconv.Atoi(c.PostForm("CategoryID"))
	title := c.PostForm("Title")
	description := c.PostForm("Description")

	// Call the use case to upload the video
	if err := u.VideoUseCase.UploadVideo(int(userID), categoryID, title, description, file); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not upload video", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully uploaded video", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

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
