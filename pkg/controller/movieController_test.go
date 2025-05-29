package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"movierental/pkg/movie/movieExternalApi"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockMovieService struct {
	ListAllMoviesFunc   func(queryParams map[string]string) ([]movieExternalApi.Movie, error)
	GetMovieDetailsFunc func(movieId string) (movieExternalApi.Movie, error)
}

func (m *MockMovieService) ListAllMovies(queryParams map[string]string) ([]movieExternalApi.Movie, error) {
	if m.ListAllMoviesFunc != nil {
		return m.ListAllMoviesFunc(queryParams)
	}
	return nil, errors.New("ListAllMoviesFunc not implemented")
}

func (m *MockMovieService) GetMovieDetails(movieId string) (movieExternalApi.Movie, error) {
	if m.GetMovieDetailsFunc != nil {
		return m.GetMovieDetailsFunc(movieId)
	}
	return movieExternalApi.Movie{}, errors.New("GetMovieDetailsFunc not implemented")
}

func setupTestRouterForMovie(mockMovieService *MockMovieService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	movieController := &MovieController{MovieService: mockMovieService}
	router.GET("/listallmovies", movieController.ListAllMovies)
	router.GET("/movie", movieController.MovieDetails)
	return router
}

func TestListAllMovies(t *testing.T) {
	t.Run("successful movie list retrieval", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			ListAllMoviesFunc: func(queryParams map[string]string) ([]movieExternalApi.Movie, error) {
				return []movieExternalApi.Movie{
					{ID: 1, Title: "Mock Movie 1", Year: 2020},
					{ID: 2, Title: "Mock Movie 2", Year: 2021},
				}, nil
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/listallmovies?limit=10&page=1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var movies []movieExternalApi.Movie
		err := json.Unmarshal(w.Body.Bytes(), &movies)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(movies) != 2 {
			t.Errorf("Expected 2 movies, got %d", len(movies))
		}
		if movies[0].Title != "Mock Movie 1" {
			t.Errorf("Expected movie title 'Mock Movie 1', got '%s'", movies[0].Title)
		}
	})

	t.Run("invalid limit parameter", func(t *testing.T) {
		mockMovieService := &MockMovieService{}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/listallmovies?limit=abc", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Invalid 'limit' parameter. Must be an integer." {
			t.Errorf("Expected error 'Invalid 'limit' parameter. Must be an integer.', got %v", response["error"])
		}
	})

	t.Run("service returns no movies found", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			ListAllMoviesFunc: func(queryParams map[string]string) ([]movieExternalApi.Movie, error) {
				return nil, errors.New("external API returned non-OK status: ok, Message: No movies were found that matched the criteria.")
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/listallmovies?genre=nonexistent", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "No movies found matching the criteria." {
			t.Errorf("Expected error 'No movies found matching the criteria.', got %v", response["error"])
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			ListAllMoviesFunc: func(queryParams map[string]string) ([]movieExternalApi.Movie, error) {
				return nil, errors.New("error calling external RapidAPI: connection refused")
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/listallmovies", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Failed to retrieve movies." {
			t.Errorf("Expected error 'Failed to retrieve movies.', got %v", response["error"])
		}
	})
}

func TestMovieDetails(t *testing.T) {
	t.Run("successful movie details retrieval", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			GetMovieDetailsFunc: func(movieId string) (movieExternalApi.Movie, error) {
				return movieExternalApi.Movie{ID: 123, Title: "Detailed Mock Movie", Year: 2022}, nil
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movie?movie_id=123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var movie movieExternalApi.Movie
		err := json.Unmarshal(w.Body.Bytes(), &movie)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if movie.ID != 123 {
			t.Errorf("Expected movie ID 123, got %d", movie.ID)
		}
		if movie.Title != "Detailed Mock Movie" {
			t.Errorf("Expected movie title 'Detailed Mock Movie', got '%s'", movie.Title)
		}
	})

	t.Run("missing movie_id parameter", func(t *testing.T) {
		mockMovieService := &MockMovieService{}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movie", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Missing required query parameter: movie_id" {
			t.Errorf("Expected error 'Missing required query parameter: movie_id', got %v", response["error"])
		}
	})

	t.Run("service returns movie not found", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			GetMovieDetailsFunc: func(movieId string) (movieExternalApi.Movie, error) {
				return movieExternalApi.Movie{}, errors.New("external API returned non-OK status: ok, Message: Movie not found!")
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movie?movie_id=999", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Movie not found!" {
			t.Errorf("Expected error 'Movie not found!', got %v", response["error"])
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockMovieService := &MockMovieService{
			GetMovieDetailsFunc: func(movieId string) (movieExternalApi.Movie, error) {
				return movieExternalApi.Movie{}, errors.New("error calling external RapidAPI: connection refused")
			},
		}
		router := setupTestRouterForMovie(mockMovieService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movie?movie_id=123", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Failed to retrieve movie details." {
			t.Errorf("Expected error 'Failed to retrieve movie details.', got %v", response["error"])
		}
	})
}
