package services

import (
	"errors"
	"fmt"
	"movierental/pkg/database"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"movierental/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct{}

func (us *UserService) CreateUser(userReq requests.CreateUser) (map[string]interface{}, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to initiate transaction for user creation: %w", tx.Error)
	}

	newUserID := uuid.New().String()
	newCartID := uuid.New().String()

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to hash password for user '%s': %w", userReq.Username, err)
	}

	userEntity := models.User{
		ID:       newUserID,
		Username: userReq.Username,
		Email:    userReq.Email,
		Password: hashedPassword,
	}

	userResult := tx.Create(&userEntity)
	if userResult.Error != nil {
		tx.Rollback()
		if errors.Is(userResult.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("username or email already exists. Please choose a different one")
		}
		return nil, fmt.Errorf("failed to create user '%s': %w", userReq.Username, userResult.Error)
	}

	cartEntity := requests.Cart{
		Id:     newCartID,
		UserId: newUserID,
		Movies: []requests.CartMovieItem{},
	}

	cartResult := tx.Create(&cartEntity)
	if cartResult.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create cart for user '%s' (ID: %s): %w", userReq.Username, newUserID, cartResult.Error)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to finalize user and cart creation: %w", err)
	}

	return gin.H{
		"message":  "User and cart created successfully!",
		"user_id":  userEntity.ID,
		"username": userEntity.Username,
		"email":    userEntity.Email,
		"cart_id":  cartEntity.Id,
	}, nil
}

func (us *UserService) LoginUser(loginReq models.User) (map[string]interface{}, error) {
	var user models.User
	err := database.DB.Where("email = ?", loginReq.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user '%s' not found. Please ensure the user exists", loginReq.Username)
		}
		return nil, fmt.Errorf("failed to retrieve user for login: %w", err)
	}

	passwordIsValid := utils.CheckPasswordHash(loginReq.Password, user.Password)
	if !passwordIsValid {
		return nil, errors.New("incorrect password")
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token for login: %w", err)
	}

	return gin.H{
		"message": "Login successful!",
		"token":   token,
	}, nil
}
