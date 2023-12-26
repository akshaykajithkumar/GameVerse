package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)
	engine.Use(middleware.AdminAuthMiddleware)
	engine.GET("/reports", adminHandler.GetReports)
	{
		usermanagement := engine.Group("/users")
		{
			usermanagement.POST("/toggle-block", adminHandler.ToggleBlockUser)

			usermanagement.GET("/getusers", adminHandler.GetUsers)
		}
		categorymanagement := engine.Group("/category")
		{
			categorymanagement.GET("/", categoryHandler.Categories)
			categorymanagement.POST("/add", categoryHandler.AddCategory)
			categorymanagement.PATCH("/update", categoryHandler.UpdateCategory)
			categorymanagement.DELETE("/delete", categoryHandler.DeleteCategory)
		}

	}
}
