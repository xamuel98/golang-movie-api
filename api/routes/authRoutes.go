package routes

import (
	"movie-api/api/resource/user/handler"

	"github.com/gin-gonic/gin"
)

// AuthRoutes creates and returns a router for handling operations on auth.
func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")

	// Define endpoints for auth
	authGroup.POST("/login", handler.LoginUser())
	authGroup.POST("/register", handler.RegisterUser())
}
