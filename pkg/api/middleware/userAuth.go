package middleware

import (
	"context"
	"fmt"
	"main/pkg/helper"
	"main/pkg/utils/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// UserAuthMiddleware is a middleware for user authentication
func UserAuthMiddleware(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	fmt.Println("Token::", token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Validate the access token
	jwttoken, err := helper.ValidateToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok || !jwttoken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}

	// Check if the access token is about to expire
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().Add(time.Minute).After(expirationTime) {
		// Access token is about to expire, refresh tokens
		refreshToken, err := c.Cookie("Refreshtoken")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Generate new access and refresh tokens
		newAccessToken, newRefreshToken, err := helper.TokensFromRefreshToken(refreshToken)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Set the new tokens in cookies
		c.SetCookie("Authorization", newAccessToken, 0, "/", "", false, true)
		c.SetCookie("Refreshtoken", newRefreshToken, 0, "/", "", false, true)
	}

	// Get user ID from claims
	userID, ok := claims["id"].(float64)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access id"})
		c.Abort()
		return
	}
	userIDString := fmt.Sprintf("%v", userID)
	var key models.UserKey = "userID"
	var val models.UserKey = models.UserKey(userIDString)

	ctx := context.WithValue(c, key, val)
	// Set the context to the request
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
