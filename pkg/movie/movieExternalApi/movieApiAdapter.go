package movieExternalApi

import (
	"log"
)

func GetAllMovies(queryParams map[string]string) []Movie {
	baseURL := "https://movie-database-api1.p.rapidapi.com"

	apiClient := NewAPIClient(baseURL)

	var moviesResponse ListMoviesResponse

	err := apiClient.Get("/list_movies.json", queryParams, &moviesResponse)
	if err != nil {
		log.Fatalf("Error fetching movies: %v", err)
	}

	if moviesResponse.Status != "ok" {
		log.Fatalf("API returned non-OK status: %s, Message: %s", moviesResponse.Status, moviesResponse.StatusMessage)
	}

	return moviesResponse.Data.Movies
}
