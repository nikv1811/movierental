package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"movierental/pkg/models/requests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockCartService struct {
	RetrieveCartFunc   func(userId interface{}) (requests.Cart, error)
	AddToCartFunc      func(userId interface{}, movieItem requests.CartMovieItem) (map[string]interface{}, error)
	RemoveFromCartFunc func(userId interface{}, movieID int) (map[string]interface{}, error)
}

func (m *MockCartService) RetrieveCart(userId interface{}) (requests.Cart, error) {
	if m.RetrieveCartFunc != nil {
		return m.RetrieveCartFunc(userId)
	}
	return requests.Cart{}, errors.New("RetrieveCartFunc not implemented")
}

func (m *MockCartService) AddToCart(userId interface{}, movieItem requests.CartMovieItem) (map[string]interface{}, error) {
	if m.AddToCartFunc != nil {
		return m.AddToCartFunc(userId, movieItem)
	}
	return nil, errors.New("AddToCartFunc not implemented")
}

func (m *MockCartService) RemoveFromCart(userId interface{}, movieID int) (map[string]interface{}, error) {
	if m.RemoveFromCartFunc != nil {
		return m.RemoveFromCartFunc(userId, movieID)
	}
	return nil, errors.New("RemoveFromCartFunc not implemented")
}

func setupTestRouterForCart(mockCartService *MockCartService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	cartController := &CartController{CartService: mockCartService}

	router.Use(func(c *gin.Context) {
		c.Set("userId", "test-user-id")
		c.Next()
	})

	router.GET("/cart", cartController.RetriveCart)
	router.POST("/cart", cartController.AddToCart)
	router.DELETE("/cart", cartController.RemoveFromCart)
	return router
}

func TestRetriveCart(t *testing.T) {
	t.Run("successful cart retrieval", func(t *testing.T) {
		mockCartService := &MockCartService{
			RetrieveCartFunc: func(userId interface{}) (requests.Cart, error) {
				return requests.Cart{
					Id:     "mock-cart-id",
					UserId: userId.(string),
					Movies: []requests.CartMovieItem{
						{ID: 101, Title: "Movie A"},
					},
				}, nil
			},
		}
		router := setupTestRouterForCart(mockCartService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var cart requests.Cart
		err := json.Unmarshal(w.Body.Bytes(), &cart)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if cart.Id != "mock-cart-id" {
			t.Errorf("Expected cart ID 'mock-cart-id', got '%s'", cart.Id)
		}
		if len(cart.Movies) != 1 {
			t.Errorf("Expected 1 movie in cart, got %d", len(cart.Movies))
		}
	})

	t.Run("internal server error during retrieval", func(t *testing.T) {
		mockCartService := &MockCartService{
			RetrieveCartFunc: func(userId interface{}) (requests.Cart, error) {
				return requests.Cart{}, errors.New("failed to retrieve cart: database error")
			},
		}
		router := setupTestRouterForCart(mockCartService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Failed to retrieve cart" {
			t.Errorf("Expected error 'Failed to retrieve cart', got %v", response["error"])
		}
	})
}

func TestAddToCart(t *testing.T) {
	t.Run("successful add to cart", func(t *testing.T) {
		mockCartService := &MockCartService{
			AddToCartFunc: func(userId interface{}, movieItem requests.CartMovieItem) (map[string]interface{}, error) {
				return map[string]interface{}{
					"message":        fmt.Sprintf("Movie '%s' (ID: %d) added to cart successfully.", movieItem.Title, movieItem.ID),
					"cart_id":        "mock-cart-id",
					"user_id":        userId.(string),
					"current_movies": []requests.CartMovieItem{movieItem},
				}, nil
			},
		}
		router := setupTestRouterForCart(mockCartService)

		movieItemJSON := `{"id":101,"title":"New Movie"}`
		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(movieItemJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["message"] != "Movie 'New Movie' (ID: 101) added to cart successfully." {
			t.Errorf("Expected message 'Movie 'New Movie' (ID: 101) added to cart successfully.', got %v", response["message"])
		}
	})

	t.Run("invalid input (missing title)", func(t *testing.T) {
		mockCartService := &MockCartService{}
		router := setupTestRouterForCart(mockCartService)

		movieItemJSON := `{"id":101}`
		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString(movieItemJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Movie ID (must be > 0) and Title are required." {
			t.Errorf("Expected error 'Movie ID (must be > 0) and Title are required.', got %v", response["error"])
		}
	})

}

func TestRemoveFromCart(t *testing.T) {
	t.Run("successful remove from cart", func(t *testing.T) {
		mockCartService := &MockCartService{
			RemoveFromCartFunc: func(userId interface{}, movieID int) (map[string]interface{}, error) {
				return map[string]interface{}{
					"message":        fmt.Sprintf("Movie with ID %d removed from cart successfully.", movieID),
					"cart_id":        "mock-cart-id",
					"user_id":        userId.(string),
					"current_movies": []requests.CartMovieItem{},
				}, nil
			},
		}
		router := setupTestRouterForCart(mockCartService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cart?movie_id=101", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["message"] != "Movie with ID 101 removed from cart successfully." {
			t.Errorf("Expected message 'Movie with ID 101 removed from cart successfully.', got %v", response["message"])
		}
	})

	t.Run("invalid movie_id parameter", func(t *testing.T) {
		mockCartService := &MockCartService{}
		router := setupTestRouterForCart(mockCartService)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cart?movie_id=abc", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "Invalid 'movie_id' parameter. Must be an integer." {
			t.Errorf("Expected error 'Invalid 'movie_id' parameter. Must be an integer.', got %v", response["error"])
		}
	})
}
