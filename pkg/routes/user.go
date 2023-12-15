package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

// func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler) {
// 	engine.POST("/login", userHandler.Login)
// 	engine.POST("/signup", userHandler.SignUp)
// 	engine.POST("/sendotp", otpHandler.SendOTP)
// 	engine.POST("/verifyotp", otpHandler.VerifyOTP)
// 	engine.POST("forgotpassword", otpHandler.ForgotPassword)

// 	// Auth middleware
// 	engine.Use(middleware.UserAuthMiddleware)
// 	engine.POST("/reportUser", userHandler.ReportUser)
// 	profile := engine.Group("/profile")
// 	{

//			profile.PATCH("/EditProfile", userHandler.EditProfile)
//			profile.GET("/GetProfile", userHandler.GetProfile)
//			settings := profile.Group("/settings")
//			profile.POST("/logout", userHandler.Logout)
//			{
//				security := settings.Group("/security")
//				{
//					security.PATCH("/change-password", userHandler.ChangePassword)
//				}
//			}
//		}
//	}
func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler) {
	engine.POST("/login", userHandler.Login)
	engine.POST("/signup", userHandler.SignUp)
	engine.POST("/sendotp", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
	engine.POST("/forgotpassword", otpHandler.ForgotPassword)

	// Auth middleware
	engine.Use(middleware.UserAuthMiddleware)

	profile := engine.Group("/profile")
	{
		profile.PATCH("/EditProfile", userHandler.EditProfile)
		profile.GET("", userHandler.GetProfile) // Change the route from "/GetProfile" to "/profile"

	}
	profile.POST("/logout", userHandler.Logout)
	engine.PATCH("/changepassword", userHandler.ChangePassword)
}
