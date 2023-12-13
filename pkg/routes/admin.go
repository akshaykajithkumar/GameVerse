package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	engine.POST("/adminlogin", adminHandler.LoginHandler)
	engine.Use(middleware.AdminAuthMiddleware)
	engine.GET("/reports", adminHandler.GetReports)
	{
		usermanagement := engine.Group("/users")
		{
			usermanagement.POST("/block", adminHandler.BlockUser)
			usermanagement.POST("/unblock", adminHandler.UnBlockUser)
			usermanagement.GET("/getusers", adminHandler.GetUsers)
		}

	}
}
