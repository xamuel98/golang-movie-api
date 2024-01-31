package routes

import (
	middleware "movie-api/api/middleware"
	"movie-api/api/resource/user/handler"

	"github.com/gin-gonic/gin"
)

// UserRoutes creates and returns a router for handling operations on user.
func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")

	// Define endpoints for user
	userGroup.Use(middleware.Authenticate())
	userGroup.GET("/:user_id", handler.GetUser())
	// userGroup.GET("/register", handler.RegisterUser())
}
