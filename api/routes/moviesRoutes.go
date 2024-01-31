package routes

import (
	middleware "movie-api/api/middleware"
	"movie-api/api/resource/movie/handler"

	"github.com/gin-gonic/gin"
)

// MoviesRoutes creates and returns a router for handling CRUD operations on movies.
func MoviesRoutes(r *gin.Engine) {
	moviesGroup := r.Group("/movies")

	// Define CRUD endpoints for movies
	moviesGroup.Use(middleware.Authenticate())
	moviesGroup.GET("/", handler.GetMovies())
	moviesGroup.GET("/:movie_id", handler.GetMovieByID())
	moviesGroup.GET("/:movie_id/cast", handler.GetMovieByIDCast())
	moviesGroup.GET("/:movie_id/similar_movies", handler.GetMovieByIDSimilarMoviesByGenre())
}
