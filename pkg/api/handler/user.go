package handler

import (
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"

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
