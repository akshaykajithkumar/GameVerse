package routes

// import (
// 	"main/pkg/api/handler"
// 	"main/pkg/api/middleware"

// 	"github.com/gin-gonic/gin"
// )

// func AdminRoutes(engine *gin.RouterGroup, adminHandler *handler.AdminHandler) {
// 	engine.POST("/adminlogin", adminHandler.LoginHandler)

// 	engine.Use(middleware.AdminAuthMiddleware)
// 	{
// 		engine.POST("/logout", adminHandler.Logout)
// 		usermanagement := engine.Group("/users")
// 		{
// 			usermanagement.POST("/block", adminHandler.BlockUser)
// 			usermanagement.POST("/unblock", adminHandler.UnBlockUser)
// 			usermanagement.GET("/getusers", adminHandler.GetUsers)
// 			usermanagement.GET("/reports", adminHandler.GetReports)

// 		}
// 	}
// }
