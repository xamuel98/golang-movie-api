package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	Movie_id          uint64           `json:"movie_id"`
	Title             string           `json:"title"`
	Overview          string           `json:"overview"`
	Popularity        float64          `json:"popularity"`
	Status            string           `json:"status"`
	Tagline           []string         `json:"tagline"`
	Video             bool             `json:"video"`
	Vote_average      float64          `json:"vote_average"`
	Vote_count        uint64           `json:"vote_count"`
	Release_date      time.Time        `json:"release_date"`
	Original_language string           `json:"original_language"`
	Spoken_languages  []SpokenLanguage `json:"spoken_languages"`
	Poster_path       string           `json:"poster_path"`
	Backdrop_path     string           `json:"backdrop_path"`
	Adult             bool             `json:"adult"`
	Genres            []Genre          `json:"genres"`
	Cast              []FullName       `json:"cast"`
	Writers           []FullName       `json:"writers"`
	Director          FullName         `json:"director"`
}

type FullName struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type Genre struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type SpokenLanguage struct {
	English_name string `json:"english_name"`
	Name         string `json:"name"`
	Iso_2        string `json:"iso_2"`
}

var movies []Movie = []Movie{
	{
		Movie_id:   1,
		Title:      "Fight Club",
		Overview:   "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy. Their concept catches on, with underground \"fight clubs\" forming in every town, until an eccentric gets in the way and ignites an out-of-control spiral toward oblivion.",
		Popularity: 61.416,
		Status:     "Released",
		Tagline: []string{
			"Mischief",
			"Mayhem",
			"Soap",
		},
		Video:             true,
		Vote_average:      8.433,
		Vote_count:        26280,
		Release_date:      time.Date(1999, 10, 15, 0, 0, 0, 0, time.UTC),
		Original_language: "en",
		Spoken_languages: []SpokenLanguage{
			{
				English_name: "English",
				Name:         "English",
				Iso_2:        "en",
			},
		},
		Poster_path:   "/path/to/poster1.jpg",
		Backdrop_path: "/path/to/backdrop1.jpg",
		Adult:         false,
		Genres: []Genre{
			{
				ID:   1,
				Name: "Drama",
			},
		},
		Cast: []FullName{
			{
				First_name: "Brad",
				Last_name:  "Pitt",
			},
			{
				First_name: "Edward",
				Last_name:  "Norton",
			},
			{
				First_name: "Meat",
				Last_name:  "Loaf",
			},
		},
		Writers: []FullName{
			{
				First_name: "Chuck",
				Last_name:  "Palahniuk",
			},
			{
				First_name: "Jim",
				Last_name:  "Uhls",
			},
		},
		Director: FullName{
			First_name: "David",
			Last_name:  "Fincher",
		},
	},
	{
		Movie_id:   2,
		Title:      "Aquaman and the Lost Kingdom",
		Overview:   "Black Manta seeks revenge on Aquaman for his father's death. Wielding the Black Trident's power, he becomes a formidable foe. To defend Atlantis, Aquaman forges an alliance with his imprisoned brother. They must protect the kingdom.",
		Popularity: 18.464,
		Status:     "Released",
		Tagline: []string{
			"Action",
			"Adventure",
			"Fantasy",
		},
		Video:             true,
		Vote_average:      5.8,
		Vote_count:        44000,
		Release_date:      time.Date(2023, 9, 24, 0, 0, 0, 0, time.UTC),
		Original_language: "en",
		Spoken_languages: []SpokenLanguage{
			{
				English_name: "English",
				Name:         "English",
				Iso_2:        "en",
			},
		},
		Poster_path:   "https://m.media-amazon.com/images/M/MV5BMzZlZDQ5NWItY2RjMC00NjRiLTlmZTgtZGE2ODEyMjVlOTJhXkEyXkFqcGdeQXVyODE5NzE3OTE@._V1_.jpg",
		Backdrop_path: "/path/to/backdrop2.jpg",
		Adult:         false,
		Genres: []Genre{
			{
				ID:   2,
				Name: "Action",
			},
			{
				ID:   3,
				Name: "Adventure",
			},
			{
				ID:   4,
				Name: "Fantasy",
			},
		},
		Cast: []FullName{
			{
				First_name: "Jason",
				Last_name:  "Momoa",
			},
			{
				First_name: "Patrick",
				Last_name:  "Wilson",
			},
			{
				First_name: "Yahya",
				Last_name:  "Abdul-Mateen II",
			},
		},
		Writers: []FullName{
			{
				First_name: "David",
				Last_name:  "Leslie Johnson-McGoldrick",
			},
			{
				First_name: "James",
				Last_name:  "Wan",
			},
			{
				First_name: "Jason",
				Last_name:  "Momoa",
			},
		},
		Director: FullName{
			First_name: "James",
			Last_name:  "Wan",
		},
	},
}

// getMovies responds with the list of all movies as JSON.
func getMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, movies)
	}
}

// Handler to get movieID
func getMovieIDHelper(c *gin.Context) (uint64, error) {
	movieIDString := c.Param("movie_id") // id is being returned as a string

	// Convert the string to uint64
	movieID, err := strconv.ParseUint(movieIDString, 10, 64)
	if err != nil {
		return 0, err
	}

	return movieID, nil
}

// Handler to get movie by ID
func getMovieByIDHelper(movieID uint64) *Movie {
	for _, movie := range movies {
		if movie.Movie_id == movieID {
			return &movie
		}
	}
	return nil
}

// Handler to find similar movies to target movie by genre
func findSimilarMoviesByGenreHelper(targetMovie *Movie, allMovies []Movie) []Movie {
	// Initialize an array to store similar movies
	var similarMovies []Movie
	// Iterate through all movies
	for _, movie := range allMovies {
		// Skip the target movie itself
		if movie.Movie_id == targetMovie.Movie_id {
			continue
		}

		// Check if the movie shares at least one genre with the target movie
		commonGenres := countCommonGenres(targetMovie.Genres, movie.Genres)

		// Adjust the criteria based on your needs (e.g., require a minimum number of shared genres)
		if commonGenres > 1 {
			similarMovies = append(similarMovies, movie)
		}
	}

	return similarMovies
}

// Helper to get common genres amongst movies
func countCommonGenres(targetGenre, currentMovieGenre []Genre) int {
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

// getMovieByID locates the movie whose ID value matches the id
// parameter sent by the client, then returns that movie as a response.
func getMovieByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the helper function to get the movie ID
		movieID, err := getMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := getMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Return movie
		c.IndentedJSON(http.StatusOK, movie)
	}
}

// getMovieByIDCast get the cast in a movie whose ID value matches the id
// parameter sent by the client, then returns that cast as a response.
func getMovieByIDCast() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the helper function to get the movie ID
		movieID, err := getMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := getMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Return movie
		c.IndentedJSON(http.StatusOK, movie.Cast)
	}
}

// getMovieByIDSimilarMoviesByGenre gets the movies whose genre matches the id
// parameter sent by the client, then returns that movies as a response.
func getMovieByIDSimilarMoviesByGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the helper function to get the movie ID
		movieID, err := getMovieIDHelper(c)
		if err != nil {
			// Handle error
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movieID format"})
			return
		}

		movie := getMovieByIDHelper(movieID)
		if movie == nil {
			// Movie was not found
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
			return
		}

		// Movie exists, find similar movies
		similarMovies := findSimilarMoviesByGenreHelper(movie, movies)

		if len(similarMovies) > 0 {
			// Return similar movies
			c.IndentedJSON(http.StatusOK, similarMovies)
		} else {
			// Return empty array movies
			c.IndentedJSON(http.StatusOK, []Movie{})
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/movies", getMovies())
	router.GET("/movies/:movie_id", getMovieByID())
	router.GET("/movies/:movie_id/cast", getMovieByIDCast())
	router.GET("/movies/:movie_id/similar_movies", getMovieByIDSimilarMoviesByGenre())

	fmt.Printf("Starting server on port: %v\n", 8080)
	router.Run() // listen and serve on 0.0.0.0:8080
}
