package handler

import (
	services "main/pkg/usecase/interface"
	"main/pkg/utils/models"
	"main/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase services.CategoryUseCase
}

func NewCategoryHandler(usecase services.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

// @Summary		Add Category
// @Description	Admin can add new categories for contents
// @Tags			Admin Content Management
// @Accept			json
// @Produce		    json
// @Param			category	query	string	true	"category"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category/add [post]
func (Cat *CategoryHandler) AddCategory(c *gin.Context) {

	cat := c.Query("category")
	CategoryResponse, err := Cat.CategoryUseCase.AddCategory(cat)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Category", CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Update Category
// @Description	Admin can update name of a category into new name
// @Tags			Admin Content Management
// @Accept			json
// @Produce		    json
// @Param			set_new_name	body	models.SetNewName	true	"set new name"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category/update [patch]
func (Cat *CategoryHandler) UpdateCategory(c *gin.Context) {

	var updateCategory models.SetNewName

	if err := c.BindJSON(&updateCategory); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	a, err := Cat.CategoryUseCase.UpdateCategory(updateCategory.Current, updateCategory.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not update the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully renamed the category", a, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		Delete Category
// @Description	Admin can delete a category
// @Tags			Admin Content Management
// @Accept			json
// @Produce		    json
// @Param			id	query	string	true	"id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category/delete [delete]
func (Cat *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID := c.Query("id")
	err := Cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary		List Categories
// @Description	Admin can view the list of  Categories
// @Tags			Admin Content Management
// @Accept			json
// @Produce		    json
// @Param			page	query  string 	true	"page"
// @Param			limit	query  string 	true	"limit"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category [get]
func (cat *CategoryHandler) Categories(c *gin.Context) {
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
	categories, err := cat.CategoryUseCase.GetCategories(page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the categories", categories, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		List Categories
// @Description	User can view the list of  Categories
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			page	query  string 	true	"page"
// @Param			limit	query  string 	true	"limit"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/users/category [get]
func (cat *CategoryHandler) CategoriesList(c *gin.Context) {
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
	categories, err := cat.CategoryUseCase.GetCategories(page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the categories", categories, nil)
	c.JSON(http.StatusOK, successRes)
}

// ListVideosByCategory is a handler for listing videos in a particular category.
// @Summary      List Videos by Category
// @Description  List videos in a specific category based on category ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        categoryID   query   int     true  "Category ID"
// @Param        page         query   string  true  "page"
// @Param        limit        query   string  true  "limit"
// @Security     Bearer
// @Success      200  {object} response.Response{}
// @Failure      400  {object} response.Response{}
// @Router       /users/category/videos [get]
func (u *CategoryHandler) ListVideosByCategory(c *gin.Context) {
	categoryIDStr := c.Query("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Invalid category ID", nil, err.Error())
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

	videos, err := u.CategoryUseCase.ListVideosByCategory(categoryID, page, limit)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not retrieve videos for the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved videos for the category", videos, nil)
	c.JSON(http.StatusOK, successRes)
}
