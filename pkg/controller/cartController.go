package controller

import (
	"fmt"
	"log"
	"strconv"

	// _ "movierental/cmd/docs"
	"movierental/pkg/database"
	"movierental/pkg/models/requests"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RetriveCart
// @Summary Retrieve user's shopping cart
// @Description Fetches the contents of the authenticated user's shopping cart.
// @Tags cart
// @Security BearerAuth // Indicates that this endpoint requires a JWT token
// @Produce application/json
// @Success 200 {object} requests.Cart "Successfully retrieved cart" // Assuming requests.Cart is the full cart structure
// @Failure 401 {object} object{error=string} "Unauthorized: unable to get userId from context or invalid token"
// @Failure 404 {object} object{error=string} "Cart not found for user"
// @Failure 500 {object} object{error=string} "Internal server error: Failed to retrieve cart"
// @Router /cart [get]
func RetriveCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to get userId from context"})
		return
	}
	var retrievedCart requests.Cart
	err := database.DB.Where("user_id = ?", userId).First(&retrievedCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Typesafe Cart not found for user_abc_123_typesafe.")
		} else {
			log.Fatalf("Failed to retrieve typesafe cart: %v", err)
		}
		return
	}

	c.JSON(200, retrievedCart)
}

// AddToCart
// @Summary Add movie to cart
// @Description Adds a specified movie item to the authenticated user's shopping cart.
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce application/json
// @Param movie_item body requests.CartMovieItem true "Movie item details to add to cart"
// @Success 200 {object} object{message=string,cart_id=string,user_id=string,current_movies=[]requests.CartMovieItem} "Movie added to cart successfully, returns updated cart details"
// @Failure 400 {object} object{error=string} "Bad Request: Invalid input or missing fields"
// @Failure 401 {object} object{error=string} "Unauthorized: unable to get userId from context or invalid token"
// @Failure 404 {object} object{error=string} "Not Found: Cart for user not found"
// @Failure 409 {object} object{error=string} "Conflict: Movie already in cart"
// @Failure 500 {object} object{error=string} "Internal server error: Database error or failed to save cart"
// @Router /cart [post]
func AddToCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to get userId from context"})
		return
	}
	var movieItem requests.CartMovieItem
	if err := c.ShouldBindJSON(&movieItem); err != nil {
		log.Printf("Validation error for AddToCart request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if movieItem.ID == 0 || movieItem.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID (must be > 0) and Title are required."})
		return
	}

	var existingCart requests.Cart
	err := database.DB.Where("user_id = ?", userId).First(&existingCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Cart not found for user ID: %s. Cannot add item to non-existent cart.", userId)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cart for user ID '%s' not found. Please ensure the user exists and their cart is created.", userId)})
		} else {
			log.Printf("Database error retrieving cart for user ID %s: %v", userId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart for adding item."})
		}
		return
	}

	for _, item := range existingCart.Movies {
		if item.ID == movieItem.ID {
			log.Printf("Movie ID %d is already in cart for user %s", movieItem.ID, userId)
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Movie '%s' (ID: %d) is already in your cart.", movieItem.Title, movieItem.ID)})
			return
		}
	}

	existingCart.Movies = append(existingCart.Movies, movieItem)

	saveResult := database.DB.Save(&existingCart)
	if saveResult.Error != nil {
		log.Printf("Database error saving updated cart for user ID %s: %v", userId, saveResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie to cart due to a database error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        fmt.Sprintf("Movie '%s' (ID: %d) added to cart successfully.", movieItem.Title, movieItem.ID),
		"cart_id":        existingCart.Id,
		"user_id":        existingCart.UserId,
		"current_movies": existingCart.Movies,
	})
}

// RemoveFromCart
// @Summary Remove movie from cart
// @Description Removes a specified movie item from the authenticated user's shopping cart by its ID.
// @Tags cart
// @Security BearerAuth
// @Produce application/json
// @Param movie_id query int true "ID of the movie to remove from cart"
// @Success 200 {object} object{message=string,cart_id=string,user_id=string,current_movies=[]requests.CartMovieItem} "Movie removed from cart successfully, returns updated cart details"
// @Failure 400 {object} object{error=string} "Bad Request: Invalid movie_id parameter"
// @Failure 401 {object} object{error=string} "Unauthorized: unable to get userId from context or invalid token"
// @Failure 404 {object} object{error=string} "Not Found: Cart for user not found or Movie not found in cart"
// @Failure 500 {object} object{error=string} "Internal server error: Database error or failed to save cart"
// @Router /cart [delete]
func RemoveFromCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to get userId from context"})
		return
	}
	// userIDStr := userId.(string)

	movieIDStr := c.Query("movie_id")
	if movieIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter: movie_id"})
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'movie_id' parameter. Must be an integer."})
		return
	}

	var existingCart requests.Cart
	err = database.DB.Where("user_id = ?", userId).First(&existingCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Cart not found for user ID: %s. Cannot remove item from non-existent cart.", userId)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cart for user ID '%s' not found.", userId)})
		} else {
			log.Printf("Database error retrieving cart for user ID %s: %v", userId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart for removing item."})
		}
		return
	}

	foundIndex := -1
	for i, item := range existingCart.Movies {
		if item.ID == movieID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Movie with ID %d not found in your cart.", movieID)})
		return
	}

	existingCart.Movies = append(existingCart.Movies[:foundIndex], existingCart.Movies[foundIndex+1:]...)

	saveResult := database.DB.Save(&existingCart)
	if saveResult.Error != nil {
		log.Printf("Database error saving updated cart for user ID %s: %v", userId, saveResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove movie from cart due to a database error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        fmt.Sprintf("Movie with ID %d removed from cart successfully.", movieID),
		"cart_id":        existingCart.Id,
		"user_id":        existingCart.UserId,
		"current_movies": existingCart.Movies,
	})
}
