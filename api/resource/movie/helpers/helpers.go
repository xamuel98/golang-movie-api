package helpers

import (
	"strconv"

	models "movie-api/api/resource/movie/model"

	"github.com/gin-gonic/gin"
)

// Handler to get movieID
func GetMovieIDHelper(c *gin.Context) (uint64, error) {
	movieIDString := c.Param("movie_id") // id is being returned as a string

	// Convert the string to uint64
	movieID, err := strconv.ParseUint(movieIDString, 10, 64)
	if err != nil {
		return 0, err
	}

	return movieID, nil
}

// Handler to get movie by ID
func GetMovieByIDHelper(movieID uint64) *models.Movie {
	for _, movie := range models.Movies {
		if movie.Movie_id == movieID {
			return &movie
		}
	}
	return nil
}

// Handler to find similar movies to target movie by genre
func FindSimilarMoviesByGenreHelper(targetMovie *models.Movie, allMovies []models.Movie) []models.Movie {
	// Initialize an array to store similar movies
	var similarMovies []models.Movie
	// Iterate through all movies
	for _, movie := range allMovies {
		// Skip the target movie itself
		if movie.Movie_id == targetMovie.Movie_id {
			continue
		}

		// Check if the movie shares at least one genre with the target movie
		commonGenres := CountCommonGenres(targetMovie.Genres, movie.Genres)

		// Adjust the criteria based on your needs (e.g., require a minimum number of shared genres)
		if commonGenres >= 1 {
			similarMovies = append(similarMovies, movie)
		}
	}

	return similarMovies
}

// Helper to get common genres amongst movies
func CountCommonGenres(targetGenre, currentMovieGenre []models.Genre) int {
	var commonGenres int = 0

	// Create a map to efficiently check the presence of genres
	genreMap := make(map[uint64]bool)
	for _, genre := range targetGenre {
		genreMap[genre.ID] = true
	}

	// Check for common genres
	for _, genre := range currentMovieGenre {
		if genreMap[genre.ID] {
			commonGenres += 1
		}
	}

	return commonGenres
}
