package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct {
	CreateUserFunc func(userReq requests.CreateUser) (map[string]interface{}, error)
	LoginUserFunc  func(loginReq models.User) (map[string]interface{}, error)
}

func (m *MockUserService) CreateUser(userReq requests.CreateUser) (map[string]interface{}, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(userReq)
	}
	return nil, errors.New("CreateUserFunc not implemented")
}

func (m *MockUserService) LoginUser(loginReq models.User) (map[string]interface{}, error) {
	if m.LoginUserFunc != nil {
		return m.LoginUserFunc(loginReq)
	}
	return nil, errors.New("LoginUserFunc not implemented")
}

func setupTestRouterForUser(mockUserService *MockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userController := &UserController{UserService: mockUserService}
	router.POST("/users", userController.CreateUser)
	router.POST("/login", userController.LoginUser)
	return router
}

func TestCreateUser(t *testing.T) {
	t.Run("successful user creation", func(t *testing.T) {
		mockUserService := &MockUserService{
			CreateUserFunc: func(userReq requests.CreateUser) (map[string]interface{}, error) {
				return map[string]interface{}{
					"message":  "User and cart created successfully!",
					"user_id":  "mock-user-id",
					"username": userReq.Username,
					"email":    userReq.Email,
					"cart_id":  "mock-cart-id",
				}, nil
			},
		}
		router := setupTestRouterForUser(mockUserService)

		userJSON := `{"username":"newuser","email":"new@example.com","password":"password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(userJSON))
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

		if response["message"] != "User and cart created successfully!" {
			t.Errorf("Expected message 'User and cart created successfully!', got %v", response["message"])
		}
		if response["user_id"] != "mock-user-id" {
			t.Errorf("Expected user_id 'mock-user-id', got %v", response["user_id"])
		}
	})

	t.Run("invalid input (bad JSON)", func(t *testing.T) {
		mockUserService := &MockUserService{}
		router := setupTestRouterForUser(mockUserService)

		userJSON := `{"username":"newuser","email":"new@example.com",}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(userJSON))
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
		if !strings.Contains(response["error"], "invalid character") {
			t.Errorf("Expected error to contain 'invalid character', got %v", response["error"])
		}
	})

	t.Run("service returns conflict error", func(t *testing.T) {
		mockUserService := &MockUserService{
			CreateUserFunc: func(userReq requests.CreateUser) (map[string]interface{}, error) {
				return nil, errors.New("username or email already exists. Please choose a different one")
			},
		}
		router := setupTestRouterForUser(mockUserService)

		userJSON := `{"username":"existinguser","email":"existing@example.com","password":"password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "username or email already exists. Please choose a different one" {
			t.Errorf("Expected error message 'username or email already exists. Please choose a different one', got %v", response["error"])
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockUserService := &MockUserService{
			CreateUserFunc: func(userReq requests.CreateUser) (map[string]interface{}, error) {
				return nil, errors.New("failed to hash password for user creation.")
			},
		}
		router := setupTestRouterForUser(mockUserService)

		userJSON := `{"username":"anyuser","email":"any@example.com","password":"password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(userJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "failed to hash password for user creation." {
			t.Errorf("Expected error message 'failed to hash password for user creation.', got %v", response["error"])
		}
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockUserService := &MockUserService{
			LoginUserFunc: func(loginReq models.User) (map[string]interface{}, error) {
				return map[string]interface{}{
					"message": "Login successful!",
					"token":   "mock-jwt-token",
				}, nil
			},
		}
		router := setupTestRouterForUser(mockUserService)

		loginJSON := `{"email":"user@example.com","password":"correctpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginJSON))
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

		if response["message"] != "Login successful!" {
			t.Errorf("Expected message 'Login successful!', got %v", response["message"])
		}
		if response["token"] != "mock-jwt-token" {
			t.Errorf("Expected token 'mock-jwt-token', got %v", response["token"])
		}
	})

	t.Run("invalid input (bad JSON)", func(t *testing.T) {
		mockUserService := &MockUserService{}
		router := setupTestRouterForUser(mockUserService)

		loginJSON := `{"email":"user@example.com",}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginJSON))
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
		if !strings.Contains(response["error"], "invalid character") {
			t.Errorf("Expected error to contain 'invalid character', got %v", response["error"])
		}
	})

	t.Run("service returns incorrect password error", func(t *testing.T) {
		mockUserService := &MockUserService{
			LoginUserFunc: func(loginReq models.User) (map[string]interface{}, error) {
				return nil, errors.New("incorrect password")
			},
		}
		router := setupTestRouterForUser(mockUserService)

		loginJSON := `{"email":"user@example.com","password":"wrongpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "incorrect password" {
			t.Errorf("Expected error message 'incorrect password', got %v", response["error"])
		}
	})

	t.Run("service returns user not found error", func(t *testing.T) {
		mockUserService := &MockUserService{
			LoginUserFunc: func(loginReq models.User) (map[string]interface{}, error) {
				return nil, errors.New("user '' not found. Please ensure the user exists")
			},
		}
		router := setupTestRouterForUser(mockUserService)

		loginJSON := `{"email":"nonexistent@example.com","password":"password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "user '' not found. Please ensure the user exists" {
			t.Errorf("Expected error message 'user '' not found. Please ensure the user exists', got %v", response["error"])
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockUserService := &MockUserService{
			LoginUserFunc: func(loginReq models.User) (map[string]interface{}, error) {
				return nil, errors.New("failed to generate token for login.")
			},
		}
		router := setupTestRouterForUser(mockUserService)

		loginJSON := `{"email":"user@example.com","password":"password123"}`
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response["error"] != "failed to generate token for login." {
			t.Errorf("Expected error message 'failed to generate token for login.', got %v", response["error"])
		}
	})
}
