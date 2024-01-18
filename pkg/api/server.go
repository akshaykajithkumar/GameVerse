package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	handler "main/pkg/api/handler"
	"main/pkg/routes"
)

// ServerHTTP represents an HTTP server for the web application.
type ServerHTTP struct {
	engine *gin.Engine // engine is the core of the Gin web framework, responsible for routing HTTP requests and handling middleware.
}

/*
NewServerHTTP creates a new instance of ServerHTTP.

Parameters:
- userHandler: A handler for user-related operations.
- otpHandler: A handler for OTP-related operations.
- adminHandler: A handler for admin-related operations.
Returns:
- *ServerHTTP: A pointer to the newly created ServerHTTP instance.
*/
func NewServerHTTP(userHandler *handler.UserHandler, otpHandler *handler.OtpHandler, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler, videoHandler *handler.VideoHandler) *ServerHTTP {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	routes.UserRoutes(engine.Group("/users"), userHandler, otpHandler, categoryHandler, videoHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, categoryHandler, videoHandler)

	return &ServerHTTP{engine: engine}
}

/*
Start starts the HTTP server and listens on port 1245.
*/
func (sh *ServerHTTP) Start() {
	sh.engine.Run(":1245")

}
