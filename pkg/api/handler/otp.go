package handler

import (
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpUseCase services.OtpUseCase
}

func NewOtpHandler(useCase services.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		otpUseCase: useCase,
	}
}

// @Summary		Send OTP
// @Description	OTP login send otp
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			otp  body  models.OTPData true	"otp-data"
// @Success		200	{object}	response.Response{}
// @Failure		400	{object}	response.Response{}
// @Router			/users/sendotp [post]
func (ot *OtpHandler) SendOTP(c *gin.Context) {

	var phone models.OTPData
	if err := c.BindJSON(&phone); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
	}

	err := ot.otpUseCase.SendOTP(phone.PhoneNumber)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Verify OTP
// @Description	OTP login verify otp
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			otp  body  models.VerifyData  true	"otp-verify"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/verifyotp [post]
func (ot *OtpHandler) VerifyOTP(c *gin.Context) {

	var code models.VerifyData
	if err := c.BindJSON(&code); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ot.otpUseCase.VerifyOTP(code)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully verified OTP", users, nil)
	// c.SetCookie("Authorization", users.AccessToken, 3600, "", "", false, false)
	c.SetCookie("Authorization", users.AccessToken, 0, "/", "", false, true)
	c.SetCookie("Refreshtoken", users.RefreshToken, 0, "/", "", false, true)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Forgot Password
// @Description	Forgot password functionality
// @Tags			User
// @Accept		json
// @Produce		json
// @Param			forgotPasswordData  body  models.ForgotPasswordData  true  "Forgot password data"
// @Success		200	{object} response.Response{}
// @Failure		400	{object} response.Response{}
// @Router			/users/forgotpassword [post]
func (ot *OtpHandler) ForgotPassword(c *gin.Context) {
	// Bind the JSON request body to the ForgotPasswordData model
	var forgotPasswordData models.ForgotPasswordData
	if err := c.BindJSON(&forgotPasswordData); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Call the ForgotPassword use case function to handle the forgot password logic
	tokenUser, err := ot.otpUseCase.ForgotPassword(forgotPasswordData)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not process forgot password request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	// Set cookies for the obtained tokens
	c.SetCookie("Authorization", tokenUser.AccessToken, 0, "/", "", false, true)
	c.SetCookie("Refreshtoken", tokenUser.RefreshToken, 0, "/", "", false, true)

	// Return a success response
	successRes := response.ClientResponse(http.StatusOK, "Forgot password request processed successfully", tokenUser, nil)
	c.JSON(http.StatusOK, successRes)
}
