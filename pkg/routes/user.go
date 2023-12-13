package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler) {
	engine.POST("/login", userHandler.Login)
	engine.POST("/signup", userHandler.SignUp)
	engine.POST("/otplogin", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)

	// Auth middleware
	engine.Use(middleware.UserAuthMiddleware)
	engine.POST("/reportUser", userHandler.ReportUser)
	profile := engine.Group("/profile")
	{
		profile.POST("/AddBio", userHandler.AddBio)
		profile.PATCH("/EditProfile", userHandler.EditProfile)
		profile.GET("/GetProfile", userHandler.GetProfile)
		settings := profile.Group("/settings")
		profile.POST("/logout", userHandler.Logout)
		{
			security := settings.Group("/security")
			{
				security.PATCH("/change-password", userHandler.ChangePassword)
			}
		}
	}
}

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
