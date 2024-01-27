package handler

import (
	"fmt"
	"main/pkg/helper"
	services "main/pkg/usecase/interface"
	"main/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	SubscriptioneUseCase services.SubscriptionUseCase
}

func NewSubscriptionHandler(usecase services.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptioneUseCase: usecase,
	}
}

// @Summary Choose Plan
// @Description User can choose a subscription plan
// @Tags User
// @Accept json
// @Produce json
// @Param creator_id query int true "creator ID"
// @Param plan_id query int true "Subscription plan ID"
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /users/plans/choose-plan [post]
func (s *SubscriptionHandler) ChoosePlan(c *gin.Context) {
	// Get user ID from the context
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println("user:", userID)
	// Get creator ID from the query parameters
	creatorIDStr := c.Query("creator_id")
	if creatorIDStr == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Creator ID is required", nil, "Creator ID is missing")
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid creator ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Get plan ID from the query parameters
	planIDStr := c.Query("plan_id")
	if planIDStr == "" {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Plan ID is required", nil, "Plan ID is missing")
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid plan ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Log parameters before calling PurchasePlan
	fmt.Println("user:", userID, creatorID, planID)

	// Add more log statements as needed for debugging

	retString, err := s.SubscriptioneUseCase.PurchasePlan(planID, creatorID, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not make the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully chose the subscription plan", retString, nil)
	c.JSON(http.StatusOK, successRes)
}

func (s *SubscriptionHandler) MakePaymentRazorPay(c *gin.Context) {

	planID := c.Query("subscription_list_id")
	fmt.Println("sublist is :", planID)
	userID, err := helper.GetUserID(c)
	fmt.Println("====", userID, planID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orderDetail, err := s.SubscriptioneUseCase.MakePaymentRazorPay(planID, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not generate order details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	c.HTML(http.StatusOK, "razorpay.html", orderDetail)
}

func (s *SubscriptionHandler) VerifyPayment(c *gin.Context) {

	planID := c.Query("order_id")
	paymentID := c.Query("payment_id")
	razorID := c.Query("razor_id")
	fmt.Println(paymentID, razorID, planID)
	err := s.SubscriptioneUseCase.VerifyPayment(paymentID, razorID, planID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not update payment details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	//clear cart
	successRes := response.ClientResponse(http.StatusOK, "Successfully updated payment details", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// // HandleUpdateSubscriptionStatus handles the manual update trigger
// func (h *SubscriptionHandler) HandleUpdateSubscriptionStatus(c *gin.Context) {
// 	// Perform any authentication or authorization checks if needed

// 	// Start the subscription update job
// 	go h.SubscriptioneUseCase.StartSubscriptionUpdateJob()

// 	// Respond with a success message
// 	c.JSON(http.StatusOK, gin.H{"message": "Subscription update job started successfully"})
// }

// HandleActivationJob handles the job triggered by the cron job

// @Summary Get Analytics
// @Description Get analytics data for subscribers count, revenue, and more
// @Tags User
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Security Bearer
// @Success 200 {object} response.Response{data=models.AnalyticsData}
// @Failure 400 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /users/analytics [get]
func (s *SubscriptionHandler) GetAnalytics(c *gin.Context) {
	userID, err := helper.GetUserID(c)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Get analytics data from the use case
	analyticsData, err := s.SubscriptioneUseCase.GetAnalytics(userID, startDate, endDate)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Error fetching analytics data", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// Create success response with analytics data
	successRes := response.Response{
		Message: "Analytics data retrieved successfully",
		Data:    analyticsData,
	}

	// Send the response
	c.JSON(http.StatusOK, successRes)
}
