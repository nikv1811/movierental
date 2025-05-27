package controller

import (
	"log"
	"movierental/pkg/movie/movieExternalApi"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	baseURL := "https://movie-database-api1.p.rapidapi.com"
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

	baseURL := "https://movie-database-api1.p.rapidapi.com"
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
