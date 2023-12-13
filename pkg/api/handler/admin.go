package handler

import (
	services "main/pkg/usecase/interface"
	models "main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(usecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

// @Summary		Admin Login
// @Description	Login handler for admins
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Param			admin	body		models.AdminLogin	true	"Admin login details"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/adminlogin [post]
func (ad *AdminHandler) LoginHandler(c *gin.Context) { // login handler for the admin

	// var adminDetails models.AdminLogin
	var adminDetails models.AdminLogin
	if err := c.BindJSON(&adminDetails); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "cannot authenticate user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", admin.Token, 3600, "/", "teeverse.online", true, false)
	c.SetCookie("Authorization", admin.AccessToken, 0, "/", "", false, true)
	c.SetCookie("Refreshtoken", admin.RefreshToken, 0, "/", "", false, true)

	successRes := response.ClientResponse(http.StatusOK, "Admin authenticated successfully", admin, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Get Users
// @Description	Retrieve users with pagination
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			limit	query		string	true	"limit"
// @Param			page	query		string	true	"Page number"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/users/getusers [get]
func (ad *AdminHandler) GetUsers(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := ad.adminUseCase.GetUsers(page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Block User
// @Description	using this handler admins can block an user
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			id	query		string	true	"user-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/users/block [post]
func (ad *AdminHandler) BlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		UnBlock an existing user
// @Description	UnBlock user
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			id	query		string	true	"user-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/users/unblock [POST]
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Get Reports List
// @Description	Get a paginated list of user reports
// @Tags			Admin
// @Accept		json
// @Produce		json
// @Param			page	query	int	false	"Page number (default: 1)"
// @Param			limit	query	int	false	"Number of items per page (default: 10)"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/admin/reports [get]
func (u *AdminHandler) GetReports(c *gin.Context) {
	page, limit := getPaginationParams(c)

	reports, err := u.adminUseCase.GetReports(page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get reports", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the list of reports
	successRes := response.ClientResponse(http.StatusOK, "Reports retrieved successfully", reports, nil)
	c.JSON(http.StatusOK, successRes)
}

func getPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	return page, limit
}
