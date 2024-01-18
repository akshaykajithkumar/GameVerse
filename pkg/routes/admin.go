package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, videoHandler *handler.VideoHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)
	engine.Use(middleware.AdminAuthMiddleware)
	engine.POST("/addtags", videoHandler.AddTagsHandler)
	engine.DELETE("/deletetags", videoHandler.DeleteTagHandler)
	engine.GET("/tags", videoHandler.GetTagsHandler)
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
