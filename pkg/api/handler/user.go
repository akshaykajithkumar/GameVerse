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
// @Router			/users/profile/settings/security/change-password [patch]
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

// @Summary		Add Bio to User
// @Description	Add a bio to the user profile
// @Tags			User
// @Accept		json
// @Produce		json
// @Param			bio	query	string	true	"Bio to be added"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/users/profile/AddBio [post]
func (u *UserHandler) AddBio(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	bio := c.Query("bio")
	if bio == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Bio cannot be empty", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := u.userUseCase.AddBio(userID, bio); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add bio to the user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully added bio to the user", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Edit User Profile
// @Description	Edit the user profile including name, email, username, phone, and bio
// @Tags	 		User
// @Accept		json
// @Produce		json
// @Param			name		query	string	false	"Updated Name"
// @Param			email		query	string	false	"Updated Email"
// @Param			username	query	string	false	"Updated Username"
// @Param			phone		query	string	false	"Updated Phone"
// @Param			bio			query	string	false	"Updated Bio"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/users/profile/EditProfile [patch]
func (u *UserHandler) EditProfile(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Extract updated fields from query parameters
	name := c.Query("name")
	email := c.Query("email")
	username := c.Query("username")
	phone := c.Query("phone")
	bio := c.Query("bio")

	// You can validate the fields as needed

	// Call the use case to update the user profile
	if err := u.userUseCase.EditProfile(userID, name, email, username, phone, bio); err != nil {
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
// @Router			/users/profile/GetProfile [get]
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
// @Router			/users/profile/logout [post]
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
