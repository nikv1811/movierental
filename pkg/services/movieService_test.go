package services

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type MockAPIClient struct {
	MockGet func(path string, queryParams map[string]string, result interface{}) error
}

func (m *MockAPIClient) Get(path string, queryParams map[string]string, result interface{}) error {
	if m.MockGet != nil {
		return m.MockGet(path, queryParams, result)
	}
	return errors.New("MockGet not implemented in test setup")
}

func TestMovieService_ListAllMovies(t *testing.T) {
	mockClient := &MockAPIClient{}

	movieService := NewMovieService(mockClient)

	t.Run("successful movie list retrieval", func(t *testing.T) {
		mockClient.MockGet = func(path string, queryParams map[string]string, result interface{}) error {
			respData := map[string]interface{}{
				"status": "ok",
				"data": map[string]interface{}{
					"movies": []map[string]interface{}{
						{"id": 1, "title": "Test Movie 1", "year": 2020},
						{"id": 2, "title": "Test Movie 2", "year": 2021},
					},
				},
			}
			jsonBytes, err := json.Marshal(respData)
			if err != nil {
				t.Fatalf("Failed to marshal mock response: %v", err)
			}
			if err := json.Unmarshal(jsonBytes, result); err != nil {
				t.Fatalf("Failed to unmarshal mock response into result: %v", err)
			}
			return nil
		}

		queryParams := map[string]string{"limit": "2", "page": "1"}
		movies, err := movieService.ListAllMovies(queryParams)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if len(movies) != 2 {
			t.Errorf("Expected 2 movies, got %d", len(movies))
		}
		if movies[0].Title != "Test Movie 1" {
			t.Errorf("Expected first movie title 'Test Movie 1', got '%s'", movies[0].Title)
		}
		if movies[1].Year != 2021 {
			t.Errorf("Expected second movie year 2021, got %d", movies[1].Year)
		}
	})

	t.Run("external API returns non-OK status", func(t *testing.T) {
		mockClient.MockGet = func(path string, queryParams map[string]string, result interface{}) error {
			respData := map[string]interface{}{
				"status":         "error",
				"status_message": "Invalid API Key",
			}
			jsonBytes, err := json.Marshal(respData)
			if err != nil {
				t.Fatalf("Failed to marshal mock response: %v", err)
			}
			if err := json.Unmarshal(jsonBytes, result); err != nil {
				t.Fatalf("Failed to unmarshal mock response into result: %v", err)
			}
			return nil
		}

		queryParams := map[string]string{"limit": "1", "page": "1"}
		_, err := movieService.ListAllMovies(queryParams)

		if err == nil {
			t.Error("Expected an error, got none")
		}
		expectedErrSubstring := "external API returned non-OK status: error, Message: Invalid API Key"
		if err != nil && !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected error to contain '%s', got: %v", expectedErrSubstring, err)
		}
	})

	t.Run("external API call fails (network error)", func(t *testing.T) {
		mockClient.MockGet = func(path string, queryParams map[string]string, result interface{}) error {
			return errors.New("network error occurred")
		}

		queryParams := map[string]string{"limit": "1", "page": "1"}
		_, err := movieService.ListAllMovies(queryParams)

		if err == nil {
			t.Error("Expected an error, got none")
		}
		expectedErrSubstring := "error calling external RapidAPI: network error occurred"
		if err != nil && !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected error to contain '%s', got: %v", expectedErrSubstring, err)
		}
	})
}

func TestMovieService_GetMovieDetails(t *testing.T) {
	mockClient := &MockAPIClient{}
	movieService := NewMovieService(mockClient)

	t.Run("successful movie details retrieval", func(t *testing.T) {
		mockClient.MockGet = func(path string, queryParams map[string]string, result interface{}) error {
			respData := map[string]interface{}{
				"status": "ok",
				"data": map[string]interface{}{
					"movie": map[string]interface{}{
						"id": 123, "title": "Detailed Movie", "year": 2022, "description_full": "A great movie.",
					},
				},
			}
			jsonBytes, err := json.Marshal(respData)
			if err != nil {
				t.Fatalf("Failed to marshal mock response: %v", err)
			}
			if err := json.Unmarshal(jsonBytes, result); err != nil {
				t.Fatalf("Failed to unmarshal mock response into result: %v", err)
			}
			return nil
		}

		movie, err := movieService.GetMovieDetails("123")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if movie.ID != 123 {
			t.Errorf("Expected movie ID 123, got %d", movie.ID)
		}
		if movie.Title != "Detailed Movie" {
			t.Errorf("Expected movie title 'Detailed Movie', got '%s'", movie.Title)
		}
	})

	t.Run("external API call fails for details (network error)", func(t *testing.T) {
		mockClient.MockGet = func(path string, queryParams map[string]string, result interface{}) error {
			return errors.New("connection refused")
		}

		_, err := movieService.GetMovieDetails("123")
		if err == nil {
			t.Error("Expected an error, got none")
		}
		expectedErrSubstring := "error calling external RapidAPI: connection refused"
		if err != nil && !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected error to contain '%s', got: %v", expectedErrSubstring, err)
		}
	})
}
