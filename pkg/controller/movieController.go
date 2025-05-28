package controller

import (
	"log"
	"movierental/config"
	"movierental/pkg/movie/movieExternalApi"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListAllMovies
// @Summary List all available movies
// @Description Retrieves a list of movies from an external API, with optional filtering and pagination.
// @Tags movies
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Number of movies to return per page (default: 20)" default(20)
// @Param page query int false "Page number for pagination (default: 1)" default(1)
// @Param minimum_rating query number false "Minimum IMDb rating for movies (e.g., 6.5)" default(0)
// @Param quality query string false "Movie quality (e.g., 720p, 1080p, 2160p, 3D)"
// @Param genre query string false "Movie genre (e.g., Action, Comedy, Horror)"
// @Param query_term query string false "Search term for movie title"
// @Param sort_by query string false "Field to sort by (e.g., title, year, rating, downloads)"
// @Param order_by query string false "Order of sorting (asc or desc)" default(desc)
// @Param with_rt_ratings query string false "Include Rotten Tomatoes ratings (true/false)"
// @Success 200 {array} movieExternalApi.Movie "Successfully retrieved list of movies"
// @Failure 400 {object} map[string]string "Bad Request: Invalid query parameters"
// @Failure 401 {object} map[string]string "Unauthorized: Missing or invalid token"
// @Failure 500 {object} map[string]string "Internal Server Error: Failed to retrieve movies from external API"
// @Failure 502 {object} map[string]string "Bad Gateway: External API returned an error"
// @Router /listallmovies [get]
func ListAllMovies(c *gin.Context) {

	QueryParams := make(map[string]string)

	limitStr := c.DefaultQuery("limit", "20")
	if _, err := strconv.Atoi(limitStr); err == nil {
		QueryParams["limit"] = limitStr
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter. Must be an integer."})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	if _, err := strconv.Atoi(pageStr); err == nil {
		QueryParams["page"] = pageStr
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter. Must be an integer."})
		return
	}

	minimumRatingStr := c.DefaultQuery("minimum_rating", "0")
	if _, err := strconv.ParseFloat(minimumRatingStr, 64); err == nil {
		QueryParams["minimum_rating"] = minimumRatingStr
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'minimum_rating' parameter. Must be a number."})
		return
	}

	if quality := c.Query("quality"); quality != "" {
		QueryParams["quality"] = quality
	}

	if genre := c.Query("genre"); genre != "" {
		QueryParams["genre"] = genre
	}

	if queryTerm := c.Query("query_term"); queryTerm != "" {
		QueryParams["query_term"] = queryTerm
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		QueryParams["sort_by"] = sortBy
	}

	if orderBy := c.Query("order_by"); orderBy != "" {
		QueryParams["order_by"] = orderBy
	}

	if withRTRatings := c.Query("with_rt_ratings"); withRTRatings != "" {
		QueryParams["with_rt_ratings"] = withRTRatings
	}

	var moviesResponse movieExternalApi.ListMoviesResponse

	baseURL := config.AppConfig.MovieAPI.BaseURL
	apiClient := movieExternalApi.NewAPIClient(baseURL)

	err := apiClient.Get("/list_movies.json", QueryParams, &moviesResponse)
	if err != nil {
		log.Printf("Error calling external RapidAPI: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies from external API"})
		return
	}

	if moviesResponse.Status != "ok" {
		log.Printf("External API returned non-OK status: %s, Message: %s", moviesResponse.Status, moviesResponse.StatusMessage)
		c.JSON(http.StatusBadGateway, gin.H{"error": "External API reported an error: " + moviesResponse.StatusMessage})
		return
	}

	c.JSON(http.StatusOK, moviesResponse.Data.Movies)
}

// MovieDetails
// @Summary Get movie details by ID
// @Description Retrieves detailed information for a specific movie by its ID. Requires authentication.
// @Tags movies
// @Security BearerAuth
// @Produce json
// @Param movie_id query int true "ID of the movie to retrieve details for"
// @Success 200 {object} movieExternalApi.Movie "Successfully retrieved movie details" // Assuming movieExternalApi.Movie is the struct for a single movie
// @Failure 400 {object} map[string]string "Bad Request: Missing or invalid movie_id parameter"
// @Failure 401 {object} map[string]string "Unauthorized: Missing or invalid token"
// @Failure 500 {object} map[string]string "Internal Server Error: Failed to retrieve movie details from external API"
// @Failure 502 {object} map[string]string "Bad Gateway: External API returned an error"
// @Router /movie [get]
func MovieDetails(c *gin.Context) {
	movieId := c.Query("movie_id")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter: movie_id"})
		return
	}

	queryParams := map[string]string{
		"movie_id": movieId,
	}
	var moviesResponse movieExternalApi.MovieResponse

	baseURL := config.AppConfig.MovieAPI.BaseURL
	apiClient := movieExternalApi.NewAPIClient(baseURL)

	err := apiClient.Get("/movie_details.json", queryParams, &moviesResponse)
	if err != nil {
		log.Printf("Error calling external RapidAPI: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies from external API"})
		return
	}

	if moviesResponse.Status != "ok" {
		log.Printf("External API returned non-OK status: %s, Message: %s", moviesResponse.Status, moviesResponse.StatusMessage)
		c.JSON(http.StatusBadGateway, gin.H{"error": "External API reported an error: " + moviesResponse.StatusMessage})
		return
	}

	c.JSON(http.StatusOK, moviesResponse.Data.Movie)
}
