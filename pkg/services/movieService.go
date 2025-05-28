package services

import (
	"fmt"
	"movierental/config"
	"movierental/pkg/movie/movieExternalApi"
)

type MovieService struct{}

func (ms *MovieService) ListAllMovies(queryParams map[string]string) ([]movieExternalApi.Movie, error) {
	var moviesResponse movieExternalApi.ListMoviesResponse

	baseURL := config.AppConfig.MovieAPI.BaseURL
	apiClient := movieExternalApi.NewAPIClient(baseURL)

	err := apiClient.Get("/list_movies.json", queryParams, &moviesResponse)
	if err != nil {
		return nil, fmt.Errorf("error calling external RapidAPI: %w", err)
	}

	if moviesResponse.Status != "ok" {
		return nil, fmt.Errorf("external API returned non-OK status: %s, Message: %s", moviesResponse.Status, moviesResponse.StatusMessage)
	}

	return moviesResponse.Data.Movies, nil
}

func (ms *MovieService) GetMovieDetails(movieId string) (movieExternalApi.Movie, error) {
	queryParams := map[string]string{
		"movie_id": movieId,
	}
	var moviesResponse movieExternalApi.MovieResponse

	baseURL := config.AppConfig.MovieAPI.BaseURL
	apiClient := movieExternalApi.NewAPIClient(baseURL)

	err := apiClient.Get("/movie_details.json", queryParams, &moviesResponse)
	if err != nil {
		return movieExternalApi.Movie{}, fmt.Errorf("error calling external RapidAPI: %w", err)
	}

	if moviesResponse.Status != "ok" {
		return movieExternalApi.Movie{}, fmt.Errorf("external API returned non-OK status: %s, Message: %s", moviesResponse.Status, moviesResponse.StatusMessage)
	}

	return moviesResponse.Data.Movie, nil
}
