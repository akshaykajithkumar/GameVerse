package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, categoyHandler *handler.CategoryHandler, videohandler *handler.VideoHandler) {
	engine.POST("/login", userHandler.Login)
	engine.POST("/signup", userHandler.SignUp)
	engine.POST("/sendotp", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
	engine.POST("/forgotpassword", otpHandler.ForgotPassword)
	engine.GET("/category", categoyHandler.CategoriesList)
	engine.GET("/category/videos", categoyHandler.ListVideosByCategory)
	// Auth middleware
	engine.Use(middleware.UserAuthMiddleware)
	engine.POST("/upload/video", videohandler.UploadVideo)

	profile := engine.Group("/profile")
	{
		profile.GET("/videos", videohandler.ListVideos)
		profile.PATCH("/videos/editVideo", videohandler.EditVideoDetails)
		profile.DELETE("/videos/delete", videohandler.DeleteVideo)
		profile.PATCH("/EditProfile", userHandler.EditProfile)
		profile.GET("", userHandler.GetProfile) // Change the route from "/GetProfile" to "/profile"

	}
	profile.POST("/logout", userHandler.Logout)
	engine.PATCH("/changepassword", userHandler.ChangePassword)
}
