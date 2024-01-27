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

// @Summary		Block or unblock User
// @Description	using this handler admins can block or unblock an user
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			id	query		string	true	"user-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/users/toggle-block [post]
func (ad *AdminHandler) ToggleBlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.ToggleBlockUser(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
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

// @Summary		Add Subscription Plan
// @Description	using this handler admins can add a new subscription plan
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			name	query		string	true	"Plan Name"
// @Param			duration	query		int	true	"Plan Duration (in days)"
// @Param			price	query		float64	true	"Plan Price"
// @Success		201	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/admin/plans/add [post]
func (ad *AdminHandler) AddSubscriptionPlan(c *gin.Context) {
	// Parse parameters
	name := c.Query("name")
	duration, err := strconv.Atoi(c.Query("duration"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "invalid duration", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	price, err := strconv.ParseFloat(c.Query("price"), 64)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "invalid price", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the function to add the subscription plan
	err = ad.adminUseCase.AddSubscriptionPlan(name, duration, price)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "plan could not be added", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added the subscription plan", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

// @Summary	Delete Subscription Plan
// @Description	using this handler admins can delete a subscription plan
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Param			id	query		int	true	"Plan ID"
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/admin/plans/delete [delete]
func (ad *AdminHandler) DeleteSubscriptionPlan(c *gin.Context) {
	planID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "invalid plan ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = ad.adminUseCase.DeleteSubscriptionPlan(planID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "plan could not be deleted", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the subscription plan", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary	Get List of Subscription Plans
// @Description	using this handler admins can get the list of subscription plans
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		Bearer
// @Success		200	{object}	[]domain.SubscriptionPlan
// @Failure		400	{object}	response.Response{}
// @Router			/admin/plans [get]
func (ad *AdminHandler) GetSubscriptionPlans(c *gin.Context) {
	plans, err := ad.adminUseCase.GetSubscriptionPlans()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not fetch subscription plans", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the list of subscription plans in the response
	c.JSON(http.StatusOK, plans)
}

// @Summary		Get Reports List
// @Description	Get a paginated list of user reports and count
// @Tags			Admin
// @Accept		json
// @Produce		json
// @Param			page	query	int	false	"Page number (default: 1)"
// @Param			limit	query	int	false	"Number of items per page (default: 10)"
// @Param			userId	query	int	false	"User ID for which reports are requested"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/admin/userReports [get]
func (u *AdminHandler) GetReportsofuser(c *gin.Context) {
	page, limit := getPaginationParams(c)
	userId, _ := strconv.Atoi(c.Query("userId"))

	reports, count, err := u.adminUseCase.GetUserReports(userId, page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get reports", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Return the list of reports and count
	successRes := response.ClientResponse(http.StatusOK, "Reports retrieved successfully", gin.H{
		"reports": reports,
		"count":   count,
	}, nil)
	c.JSON(http.StatusOK, successRes)
}
