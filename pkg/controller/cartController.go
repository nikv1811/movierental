package controller

import (
	"fmt"
	"log"
	"movierental/pkg/models/requests"
	"movierental/pkg/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	CartService services.CartServiceInterface
}

// RetriveCart
// @Summary Retrieve user's shopping cart
// @Description Fetches the contents of the authenticated user's shopping cart.
// @Tags cart
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} requests.Cart "Successfully retrieved cart"
// @Failure 401 {object} object{error=string} "Unauthorized: unable to get userId from context or invalid token"
// @Failure 404 {object} object{error=string} "Cart not found for user"
// @Failure 500 {object} object{error=string} "Internal server error: Failed to retrieve cart"
// @Router /cart [get]
func (cc *CartController) RetriveCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to get userId from context"})
		return
	}

	retrievedCart, err := cc.CartService.RetrieveCart(userId)
	if err != nil {
		log.Printf("Error retrieving cart: %v", err)
		if err.Error() == "cart not found for user" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found for user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		}
		return
	}
	c.JSON(http.StatusOK, retrievedCart)
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
func (cc *CartController) AddToCart(c *gin.Context) {
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

	response, err := cc.CartService.AddToCart(userId, movieItem)
	if err != nil {
		log.Printf("Error adding to cart: %v", err)
		if err.Error() == fmt.Sprintf("cart for user ID '%s' not found. Please ensure the user exists and their cart is created", userId) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == fmt.Sprintf("movie '%s' (ID: %d) is already in your cart", movieItem.Title, movieItem.ID) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response)
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
func (cc *CartController) RemoveFromCart(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to get userId from context"})
		return
	}

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

	response, err := cc.CartService.RemoveFromCart(userId, movieID)
	if err != nil {
		log.Printf("Error removing from cart: %v", err)
		if err.Error() == fmt.Sprintf("cart for user ID '%s' not found", userId) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == fmt.Sprintf("movie with ID %d not found in your cart", movieID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response)
}
