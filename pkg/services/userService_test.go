package services

import (
	"movierental/pkg/database"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"movierental/pkg/utils"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	testDB := utils.SetupTestDB()
	defer utils.ClearTestDB(testDB)

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	userService := &UserService{}

	t.Run("Successful user creation", func(t *testing.T) {
		userReq := requests.CreateUser{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		response, err := userService.CreateUser(userReq)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if response == nil {
			t.Error("Expected a non-nil response, got nil")
		}
		if response["message"] != "User and cart created successfully!" {
			t.Errorf("Expected message 'User and cart created successfully!', got: %v", response["message"])
		}
		if userID, ok := response["user_id"].(string); !ok || userID == "" {
			t.Errorf("Expected non-empty user_id, got: %v", response["user_id"])
		}
		if cartID, ok := response["cart_id"].(string); !ok || cartID == "" {
			t.Errorf("Expected non-empty cart_id, got: %v", response["cart_id"])
		}

		var user models.User
		err = testDB.Where("email = ?", userReq.Email).First(&user).Error
		if err != nil {
			t.Errorf("Expected user to exist in DB, got error: %v", err)
		}
		if user.Username != userReq.Username {
			t.Errorf("Expected username %s, got %s", userReq.Username, user.Username)
		}

		var cart requests.Cart
		err = testDB.Where("user_id = ?", user.ID).First(&cart).Error
		if err != nil {
			t.Errorf("Expected cart to exist in DB, got error: %v", err)
		}
	})
}

func TestUserService_LoginUser(t *testing.T) {
	testDB := utils.SetupTestDB()
	defer utils.ClearTestDB(testDB)

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	userService := &UserService{}

	hashedPassword, _ := utils.HashPassword("correctpassword")
	testUser := models.User{
		ID:       "user-123",
		Username: "loginuser",
		Email:    "login@example.com",
		Password: hashedPassword,
	}
	testDB.Create(&testUser)

	t.Run("Successful login", func(t *testing.T) {
		loginReqSuccess := models.User{
			Email:    "login@example.com",
			Password: "correctpassword",
		}
		response, err := userService.LoginUser(loginReqSuccess)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if response == nil {
			t.Error("Expected a non-nil response, got nil")
		}
		if response["message"] != "Login successful!" {
			t.Errorf("Expected message 'Login successful!', got: %v", response["message"])
		}
		if token, ok := response["token"].(string); !ok || token == "" {
			t.Errorf("Expected non-empty token, got: %v", response["token"])
		}
	})
}
