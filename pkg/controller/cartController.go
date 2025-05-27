package controller

import (
	"fmt"
	"log"
	"movierental/pkg/database"
	"movierental/pkg/models/requests"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RetriveCart(c *gin.Context) {
	userId := c.Param("user_id")
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

func AddToCart(c *gin.Context) {
	userID := c.Param("user_id")

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
	err := database.DB.Where("user_id = ?", userID).First(&existingCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Cart not found for user ID: %s. Cannot add item to non-existent cart.", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cart for user ID '%s' not found. Please ensure the user exists and their cart is created.", userID)})
		} else {
			log.Printf("Database error retrieving cart for user ID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart for adding item."})
		}
		return
	}

	for _, item := range existingCart.Movies {
		if item.ID == movieItem.ID {
			log.Printf("Movie ID %d is already in cart for user %s", movieItem.ID, userID)
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Movie '%s' (ID: %d) is already in your cart.", movieItem.Title, movieItem.ID)})
			return
		}
	}

	existingCart.Movies = append(existingCart.Movies, movieItem)

	saveResult := database.DB.Save(&existingCart)
	if saveResult.Error != nil {
		log.Printf("Database error saving updated cart for user ID %s: %v", userID, saveResult.Error)
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
