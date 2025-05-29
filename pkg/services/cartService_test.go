package services

import (
	"movierental/pkg/database"
	"movierental/pkg/models/requests"
	"movierental/pkg/utils"
	"testing"
)

func TestCartService_RetrieveCart(t *testing.T) {
	testDB := utils.SetupTestDB()
	defer utils.ClearTestDB(testDB)

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	cartService := &CartService{}

	testUserID := "user-cart-1"
	utils.CreateTestUserAndCart(testDB, testUserID, "cart1@example.com")

	t.Run("Successfully retrieve cart", func(t *testing.T) {
		cart, err := cartService.RetrieveCart(testUserID)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if cart.UserId != testUserID {
			t.Errorf("Expected UserID %s, got %s", testUserID, cart.UserId)
		}
		if len(cart.Movies) != 0 {
			t.Errorf("Expected cart to be empty, got %d movies", len(cart.Movies))
		}
	})

	t.Run("Cart not found for non-existent user", func(t *testing.T) {
		_, err := cartService.RetrieveCart("nonexistent-user")
		if err == nil {
			t.Error("Expected an error for non-existent cart, got none")
		}
		expectedErr := "cart not found for user"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got: '%v'", expectedErr, err)
		}
	})
}

func TestCartService_AddToCart(t *testing.T) {
	testDB := utils.SetupTestDB()
	defer utils.ClearTestDB(testDB)

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	cartService := &CartService{}

	testUserID := "user-cart-2"
	utils.CreateTestUserAndCart(testDB, testUserID, "cart2@example.com")

	movieItem1 := requests.CartMovieItem{ID: 101, Title: "Movie A"}
	movieItem2 := requests.CartMovieItem{ID: 102, Title: "Movie B"}

	t.Run("Add first movie successfully", func(t *testing.T) {
		response, err := cartService.AddToCart(testUserID, movieItem1)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		expectedMsg := "Movie 'Movie A' (ID: 101) added to cart successfully."
		if response["message"] != expectedMsg {
			t.Errorf("Expected message '%s', got: %v", expectedMsg, response["message"])
		}
		if currentMovies, ok := response["current_movies"].([]requests.CartMovieItem); !ok || len(currentMovies) != 1 {
			t.Errorf("Expected 1 movie in cart, got %d", len(currentMovies))
		}

		var updatedCart requests.Cart
		testDB.Where("user_id = ?", testUserID).First(&updatedCart)
		if len(updatedCart.Movies) != 1 {
			t.Errorf("DB: Expected 1 movie, got %d", len(updatedCart.Movies))
		}
		if updatedCart.Movies[0].ID != movieItem1.ID {
			t.Errorf("DB: Expected movie ID %d, got %d", movieItem1.ID, updatedCart.Movies[0].ID)
		}
	})

	t.Run("Add second movie successfully", func(t *testing.T) {
		response, err := cartService.AddToCart(testUserID, movieItem2)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		expectedMsg := "Movie 'Movie B' (ID: 102) added to cart successfully."
		if response["message"] != expectedMsg {
			t.Errorf("Expected message '%s', got: %v", expectedMsg, response["message"])
		}
		if currentMovies, ok := response["current_movies"].([]requests.CartMovieItem); !ok || len(currentMovies) != 2 {
			t.Errorf("Expected 2 movies in cart, got %d", len(currentMovies))
		}

		var updatedCart requests.Cart
		testDB.Where("user_id = ?", testUserID).First(&updatedCart)
		if len(updatedCart.Movies) != 2 {
			t.Errorf("DB: Expected 2 movies, got %d", len(updatedCart.Movies))
		}
		if updatedCart.Movies[1].ID != movieItem2.ID {
			t.Errorf("DB: Expected movie ID %d, got %d", movieItem2.ID, updatedCart.Movies[1].ID)
		}
	})

	t.Run("Add duplicate movie", func(t *testing.T) {
		response, err := cartService.AddToCart(testUserID, movieItem1)
		if err == nil {
			t.Error("Expected an error for duplicate movie, got none")
		}
		if response != nil {
			t.Errorf("Expected nil response, got: %v", response)
		}

		expectedErr := "movie 'Movie A' (ID: 101) is already in your cart"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got: '%v'", expectedErr, err)
		}
	})

	t.Run("Cart not found (non-existent user)", func(t *testing.T) {
		nonExistentUserID := "nonexistent-user-id"
		response, err := cartService.AddToCart(nonExistentUserID, movieItem1)
		if err == nil {
			t.Error("Expected an error for non-existent cart, got none")
		}
		if response != nil {
			t.Errorf("Expected nil response, got: %v", response)
		}

		expectedErr := "cart for user ID 'nonexistent-user-id' not found. Please ensure the user exists and their cart is created"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got: '%v'", expectedErr, err)
		}
	})
}

func TestCartService_RemoveFromCart(t *testing.T) {
	testDB := utils.SetupTestDB()
	defer utils.ClearTestDB(testDB)

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	cartService := &CartService{}

	testUserID := "user-cart-3"
	utils.CreateTestUserAndCart(testDB, testUserID, "cart3@example.com")

	movieItem1 := requests.CartMovieItem{ID: 201, Title: "Movie X"}
	movieItem2 := requests.CartMovieItem{ID: 202, Title: "Movie Y"}

	_, _ = cartService.AddToCart(testUserID, movieItem1)
	_, _ = cartService.AddToCart(testUserID, movieItem2)

	t.Run("Remove existing movie successfully", func(t *testing.T) {
		response, err := cartService.RemoveFromCart(testUserID, movieItem1.ID)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		expectedMsg := "Movie with ID 201 removed from cart successfully."
		if response["message"] != expectedMsg {
			t.Errorf("Expected message '%s', got: %v", expectedMsg, response["message"])
		}
		if currentMovies, ok := response["current_movies"].([]requests.CartMovieItem); !ok || len(currentMovies) != 1 {
			t.Errorf("Expected 1 movie in cart, got %d", len(currentMovies))
		}

		var updatedCart requests.Cart
		testDB.Where("user_id = ?", testUserID).First(&updatedCart)
		if len(updatedCart.Movies) != 1 {
			t.Errorf("DB: Expected 1 movie, got %d", len(updatedCart.Movies))
		}
		if updatedCart.Movies[0].ID != movieItem2.ID {
			t.Errorf("DB: Expected movie ID %d, got %d", movieItem2.ID, updatedCart.Movies[0].ID)
		}
	})

	t.Run("Try to remove non-existent movie", func(t *testing.T) {
		response, err := cartService.RemoveFromCart(testUserID, 999)
		if err == nil {
			t.Error("Expected an error for non-existent movie in cart, got none")
		}
		if response != nil {
			t.Errorf("Expected nil response, got: %v", response)
		}
		expectedErr := "movie with ID 999 not found in your cart"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got: '%v'", expectedErr, err)
		}
	})

	t.Run("Remove the last movie", func(t *testing.T) {
		response, err := cartService.RemoveFromCart(testUserID, movieItem2.ID)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		expectedMsg := "Movie with ID 202 removed from cart successfully."
		if response["message"] != expectedMsg {
			t.Errorf("Expected message '%s', got: %v", expectedMsg, response["message"])
		}
		if currentMovies, ok := response["current_movies"].([]requests.CartMovieItem); !ok || len(currentMovies) != 0 {
			t.Errorf("Expected 0 movies in cart, got %d", len(currentMovies))
		}

		var updatedCart requests.Cart
		testDB.Where("user_id = ?", testUserID).First(&updatedCart)
		if len(updatedCart.Movies) != 0 {
			t.Errorf("DB: Expected 0 movies, got %d", len(updatedCart.Movies))
		}
	})

	t.Run("Cart not found (non-existent user for removal)", func(t *testing.T) {
		nonExistentUserID := "nonexistent-user-id-remove"
		response, err := cartService.RemoveFromCart(nonExistentUserID, movieItem1.ID)
		if err == nil {
			t.Error("Expected an error for non-existent cart, got none")
		}
		if response != nil {
			t.Errorf("Expected nil response, got: %v", response)
		}
		expectedErr := "cart for user ID 'nonexistent-user-id-remove' not found"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got: '%v'", expectedErr, err)
		}
	})
}
