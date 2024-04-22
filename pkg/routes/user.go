package routes

import (
	"main/pkg/api/handler"
	"main/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup, userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, categoyHandler *handler.CategoryHandler, videohandler *handler.VideoHandler, subscriptionhandler *handler.SubscriptionHandler) {
	engine.POST("/login", userHandler.Login)
	engine.POST("/signup", userHandler.SignUp)
	engine.POST("/logout", userHandler.Logout)
	engine.POST("/sendotp", otpHandler.SendOTP)
	engine.POST("/verifyotp", otpHandler.VerifyOTP)
	engine.POST("/forgotpassword", otpHandler.ForgotPassword)
	engine.GET("/category", categoyHandler.CategoriesList)
	engine.GET("/plans", userHandler.GetSubscriptionPlans)
	engine.GET("/category/videos", categoyHandler.ListVideosByCategory)
	engine.GET("/videos", videohandler.ListtVideos)
	// Auth middleware
	engine.Use(middleware.UserAuthMiddleware)
	engine.GET("/followingList", userHandler.GetFollowingList)
	engine.GET("/followersList", userHandler.GetFollowersList)
	engine.GET("/search", userHandler.SearchUsers)
	engine.POST("search/toggleFollow", userHandler.ToggleFollow)
	engine.GET("/analytics", subscriptionhandler.GetAnalytics)
	engine.POST("/reportUser", userHandler.ReportUser)
	engine.GET("/tags", videohandler.GetTagsForUserHandler)
	engine.POST("/selectTags", videohandler.StoreUserTags)
	// engine.POST("/upload/video", videohandler.UploadVideo)
	// payment := engine.Group("users/plans")

	engine.POST("plans/choose-plan", subscriptionhandler.ChoosePlan)
	engine.GET("plans/choose-plan/razorpay", subscriptionhandler.MakePaymentRazorPay)
	engine.GET("plans/update_status", subscriptionhandler.VerifyPayment)

	profile := engine.Group("/profile")
	{
		profile.GET("/videos", videohandler.ListVideos)
		profile.GET("/videos/recommendation", videohandler.RecommendationList)

		profile.GET("/videos/comments", videohandler.GetCommentsHandler)
		profile.POST("/videos/comment", videohandler.CommentVideoHandler)
		// profile.GET("/videos/watch", videohandler.WatchVideo)
		profile.PATCH("/videos/editVideo", videohandler.EditVideoDetails)
		profile.DELETE("/videos/delete", videohandler.DeleteVideo)
		profile.POST("/videos/like", videohandler.ToggleLikeVideo)
		profile.PATCH("/EditProfile", userHandler.EditProfile)
		profile.GET("", userHandler.GetProfile) // Change the route from "/GetProfile" to "/profile"

	}

	engine.PATCH("/changepassword", userHandler.ChangePassword)
}
