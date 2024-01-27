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

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

// Signup is a handler for user Registration
// @Summary		User Signup
// @Description	user can signup by giving their details
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			signup  body  models.UserDetails  true	"signup"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/signup [post]
func (i *UserHandler) SignUp(c *gin.Context) {

	var user models.UserDetails
	// bind the user details to the struct
	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// business logic goes inside this function
	userCreated, err := i.userUseCase.SignUp(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not signed up", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully signed up", userCreated, nil)
	c.JSON(http.StatusCreated, successRes)

}

// Login is a handler for user login
// @Summary		User Login
// @Description	user can log in by giving their details
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			login  body  models.UserLogin  true	"login"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/login [post]
func (i *UserHandler) Login(c *gin.Context) {
	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userDetails, err := i.userUseCase.Login(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not be logged in", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User successfully logged in", userDetails, nil)

	c.SetCookie("Authorization", userDetails.AccessToken, 0, "/", "", false, true)
	c.SetCookie("Refreshtoken", userDetails.RefreshToken, 0, "/", "", false, true)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Change Password
// @Description	user can change their password
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			changepassword  body  models.ChangePassword  true	"changepassword"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/change-password [patch]
func (i *UserHandler) ChangePassword(c *gin.Context) {

	id, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var ChangePassword models.ChangePassword
	if err := c.BindJSON(&ChangePassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.ChangePassword(id, ChangePassword.Oldpassword, ChangePassword.Password, ChangePassword.Repassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "password changed Successfully ", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// EditProfile is a handler for editing user profile details.
// @Summary      Edit User Profile
// @Description  Edit the user profile including name, email, username, phone, bio, and profile picture
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param			userProfileRequest	body	models.EditUserProfileResponse	true	"User Profile Request"
// @Param        ProfilePicture      formData  file    true  "Profile Picture"
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/profile/EditProfile [patch]
func (u *UserHandler) EditProfile(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Bind request body to UserProfileResponse struct
	var userProfileRequest models.UserProfileResponse
	if err := c.ShouldBind(&userProfileRequest); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Retrieve file from form
	file, err := c.FormFile("ProfilePicture")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Error retrieving image from form", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to update the user profile
	if err := u.userUseCase.EditProfile(userID, userProfileRequest.Name, userProfileRequest.Email, userProfileRequest.Username, userProfileRequest.Phone, userProfileRequest.Bio, file); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not edit user profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully edited user profile", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Get User Profile
// @Description	Get the user profile details
// @Tags	 		User
// @Accept		json
// @Produce		json
// @Security		Bearer
// @Success		200	{object}	models.UserProfileResponse
// @Failure		400	{object}	response.Response{}
// @Router			/users/profile [get]
func (u *UserHandler) GetProfile(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the use case to get the user profile
	userProfile, err := u.userUseCase.GetProfile(userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get user profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the user profile details in the response
	c.JSON(http.StatusOK, userProfile)
}

// Logout is a handler for user logout
// @Summary		User Logout
// @Description	Logout the currently authenticated user
// @Tags			User
// @Accept		json
// @Produce		json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/logout [post]
func (i *UserHandler) Logout(c *gin.Context) {

	// Clear the access token and refresh token cookies
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.SetCookie("Refreshtoken", "", -1, "/", "", false, true)

	successRes := response.ClientResponse(http.StatusOK, "User successfully logged out", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Report User
// @Description	Submit a report for a user
// @Tags			User
// @Accept		json
// @Produce		json
// @Param			targetUserID	query	int	true	"ID of the user being reported"
// @Param			reason		query	string	true	"Reason for the report"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/users/reportUser [post]
func (u *UserHandler) ReportUser(c *gin.Context) {
	// Get the reporter's user ID from the token
	reporterUserID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get reporter's userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get the target user's ID from the request
	targetUserID, err := strconv.Atoi(c.Query("targetUserID"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid targetUserID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get the reason for the report from the request
	reason := c.Query("reason")
	if reason == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Reason cannot be empty", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the user use case to handle the reporting logic
	if err := u.userUseCase.ReportUser(reporterUserID, targetUserID, reason); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not submit the report", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Report submitted successfully
	successRes := response.ClientResponse(http.StatusOK, "Successfully submitted the report", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Toggle Follow
// @Description	Toggle the follow status between two users
// @Tags			User
// @Accept		json
// @Produce		json
// @Param			followingUserID	query	int	true	"ID of the user being followed/unfollowed"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/users/search/toggleFollow [post]
func (u *UserHandler) ToggleFollow(c *gin.Context) {
	// Get the follower's user ID from the token
	followerID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get follower's userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get the following user's ID from the request
	followingID, err := strconv.Atoi(c.Query("followingUserID"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid followingUserID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the user use case to handle the follow/unfollow logic
	if err := u.userUseCase.ToggleFollow(followerID, followingID); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not toggle follow status", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Toggle successful
	successRes := response.ClientResponse(http.StatusOK, "Toggle follow status successful", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary	Get Following List with Pagination
// @Description	Get the paginated list of users (ID and username) that the logged-in user is following
// @Tags			User
// @Accept		json
// @Produce		json
// @Security		Bearer
// @Param			page	query	int	true	"Page number"
// @Param			limit	query	int	true	"Limit per page"
// @Success		200		{object}	[]models.FollowingUser
// @Failure		400		{object}	response.Response{}
// @Router			/users/followingList [get]
func (u *UserHandler) GetFollowingList(c *gin.Context) {
	// Get the user ID from the token
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get user ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get the page and limit from the query parameters
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Page number not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Limit not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the user use case to get the following list with pagination
	followingList, err := u.userUseCase.GetFollowingListWithPagination(userID, page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get following list", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the following list
	c.JSON(http.StatusOK, followingList)
}

// @Summary	Search Users by Name
// @Description	Search for users by name and return the results in alphabetical order with pagination
// @Tags		User
// @Accept	json
// @Produce	json
// @Security	Bearer
// @Param		searchTerm	query	string	true	"Search term for user name"
// @Param		page	query	int	true	"Page number"
// @Param		limit	query	int	true	"Limit per page"
// @Success	200	{object}	[]map[string]interface{}
// @Failure	400	{object}	response.Response{}
// @Router	/users/search [get]
func (u *UserHandler) SearchUsers(c *gin.Context) {
	// Get the search term, page, and limit from the query parameters
	searchTerm := c.Query("searchTerm")
	if searchTerm == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Search term cannot be empty", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Page number not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Limit not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the user use case to search for users by name with pagination
	searchResults, err := u.userUseCase.SearchUsersByNameWithPagination(searchTerm, page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not search for users", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	simplifiedResults := make([]map[string]interface{}, len(searchResults))
	for i, user := range searchResults {
		simplifiedResults[i] = map[string]interface{}{
			"ID":       user.ID,
			"username": user.Username,
		}
	}
	// Return the search results
	c.JSON(http.StatusOK, simplifiedResults)
}

// @Summary	Get List of Subscription Plans
// @Description	using this handler users can get the list of subscription plans
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Success		200	{object}	[]domain.SubscriptionPlan
// @Failure		400	{object}	response.Response{}
// @Router			/users/plans [get]
func (ad *UserHandler) GetSubscriptionPlans(c *gin.Context) {
	plans, err := ad.userUseCase.GetSubscriptionPlans()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not fetch subscription plans", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the list of subscription plans in the response
	c.JSON(http.StatusOK, plans)
}

// GetFollowersList retrieves the paginated list of users (ID and username) that the logged-in user's followers
// @Summary	Get Followers List with Pagination
// @Description	Get the paginated list of users (ID and username) that the logged-in user's followers
// @Tags			User
// @Accept		json
// @Produce		json
// @Security		Bearer
// @Param			page	query	int	true	"Page number"
// @Param			limit	query	int	true	"Limit per page"
// @Success		200		{object}	[]models.FollowerUser
// @Failure		400		{object}	response.Response{}
// @Router			/users/followersList [get]
func (u *UserHandler) GetFollowersList(c *gin.Context) {
	// Get the user ID from the token
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get user ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get the page and limit from the query parameters
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Page number not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Limit not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the user use case to get the followers list with pagination
	followersList, err := u.userUseCase.GetFollowersListWithPagination(userID, page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get followers list", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the followers list
	c.JSON(http.StatusOK, followersList)
}
