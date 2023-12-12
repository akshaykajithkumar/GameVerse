package routes

import (
	"main/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler /*, otpHandler *handler.OtpHandler*/) {
	// engine.POST("/login", userHandler.Login)
	engine.POST("/signup", userHandler.SignUp)
	// engine.POST("/otplogin", otpHandler.SendOTP)
	// engine.POST("/verifyotp", otpHandler.VerifyOTP)
	// engine.Use(middleware.UserAuthMiddleware)
	// {
	// 	engine.POST("/logout", userHandler.Logout)
	// 	profile := engine.Group("/profile")
	// 	{
	// 		security := profile.Group("/security")
	// 		{
	// 			security.PATCH("/change-password", userHandler.ChangePassword)
	// 		}

	// 		report := profile.Group("/report")
	// 		{
	// 			report.POST("/report-user", userHandler.ReportUser)
	// 		}

	// 	}
	// }
}
