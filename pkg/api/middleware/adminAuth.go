package middleware

import (
	"fmt"
	"main/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthCustomClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

/*
AdminAuthMiddleware is a middleware for admin authentication

Parameters:
- c: Gin Context.
*/
func AdminAuthMiddleware(c *gin.Context) {
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
	if !ok || role != "admin" {
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

	c.Next()
}
