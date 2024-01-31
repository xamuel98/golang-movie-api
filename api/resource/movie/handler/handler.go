package handler

import (
	"net/http"

	helper "movie-api/api/resource/movie/helpers"
	models "movie-api/api/resource/movie/model"

	"github.com/gin-gonic/gin"
)

// GetMovies responds with the list of all movies as JSON.
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Content-Type header to application/json
		c.Header("Content-Type", "application/json")

		c.IndentedJSON(http.StatusOK, models.Movies)
	}
}

// getMGetMovieByIDovieByID locates the movie whose ID value matches the id
// parameter sent by the client, then returns that movie as a response.
func GetMovieByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Content-Type header to application/json
		c.Header("Content-Type", "application/json")

		// Use the helper function to get the movie ID
		movieID, err := helper.GetMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := helper.GetMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Return movie
		c.IndentedJSON(http.StatusOK, movie)
	}
}

// GetMovieByIDCast get the cast in a movie whose ID value matches the id
// parameter sent by the client, then returns that cast as a response.
func GetMovieByIDCast() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Content-Type header to application/json
		c.Header("Content-Type", "application/json")

		// Use the helper function to get the movie ID
		movieID, err := helper.GetMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := helper.GetMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Return movie
		c.IndentedJSON(http.StatusOK, movie.Cast)
	}
}

// GetMovieByIDSimilarMoviesByGenre gets the movies whose genre matches the id
// parameter sent by the client, then returns that movies as a response.
func GetMovieByIDSimilarMoviesByGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Content-Type header to application/json
		c.Header("Content-Type", "application/json")

		// Use the helper function to get the movie ID
		movieID, err := helper.GetMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := helper.GetMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Movie exists, find similar movies
		similarMovies := helper.FindSimilarMoviesByGenreHelper(movie, models.Movies)

		if len(similarMovies) > 0 {
			// Return similar movies
			c.IndentedJSON(http.StatusOK, similarMovies)
		} else {
			// Return empty array movies
			c.IndentedJSON(http.StatusOK, []models.Movie{})
		}
	}
}
