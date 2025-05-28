package services

import (
	"errors"
	"fmt"
	"movierental/pkg/database"
	"movierental/pkg/models/requests"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartService struct{}

func (cs *CartService) RetrieveCart(userId interface{}) (requests.Cart, error) {
	var retrievedCart requests.Cart
	err := database.DB.Where("user_id = ?", userId).First(&retrievedCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return requests.Cart{}, errors.New("cart not found for user")
		}
		return requests.Cart{}, fmt.Errorf("failed to retrieve cart: %w", err)
	}
	return retrievedCart, nil
}

func (cs *CartService) AddToCart(userId interface{}, movieItem requests.CartMovieItem) (map[string]interface{}, error) {
	var existingCart requests.Cart
	err := database.DB.Where("user_id = ?", userId).First(&existingCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cart for user ID '%s' not found. Please ensure the user exists and their cart is created", userId)
		}
		return nil, fmt.Errorf("failed to retrieve cart for adding item: %w", err)
	}

	for _, item := range existingCart.Movies {
		if item.ID == movieItem.ID {
			return nil, fmt.Errorf("movie '%s' (ID: %d) is already in your cart", movieItem.Title, movieItem.ID)
		}
	}

	existingCart.Movies = append(existingCart.Movies, movieItem)

	saveResult := database.DB.Save(&existingCart)
	if saveResult.Error != nil {
		return nil, fmt.Errorf("failed to add movie to cart due to a database error: %w", saveResult.Error)
	}

	return gin.H{
		"message":        fmt.Sprintf("Movie '%s' (ID: %d) added to cart successfully.", movieItem.Title, movieItem.ID),
		"cart_id":        existingCart.Id,
		"user_id":        existingCart.UserId,
		"current_movies": existingCart.Movies,
	}, nil
}

func (cs *CartService) RemoveFromCart(userId interface{}, movieID int) (map[string]interface{}, error) {
	var existingCart requests.Cart
	err := database.DB.Where("user_id = ?", userId).First(&existingCart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cart for user ID '%s' not found", userId)
		}
		return nil, fmt.Errorf("failed to retrieve cart for removing item: %w", err)
	}

	foundIndex := -1
	for i, item := range existingCart.Movies {
		if item.ID == movieID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return nil, fmt.Errorf("movie with ID %d not found in your cart", movieID)
	}

	existingCart.Movies = append(existingCart.Movies[:foundIndex], existingCart.Movies[foundIndex+1:]...)

	saveResult := database.DB.Save(&existingCart)
	if saveResult.Error != nil {
		return nil, fmt.Errorf("failed to remove movie from cart due to a database error: %w", saveResult.Error)
	}

	return gin.H{
		"message":        fmt.Sprintf("Movie with ID %d removed from cart successfully.", movieID),
		"cart_id":        existingCart.Id,
		"user_id":        existingCart.UserId,
		"current_movies": existingCart.Movies,
	}, nil
}
