package controller

import (
	"errors"
	"fmt"
	"log"
	"movierental/pkg/database"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"movierental/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateUser
// @Summary Register a new user
// @Description Creates a new user account with a username, email, and password. Also creates an associated shopping cart.
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.CreateUser true "User registration details (username, email, password)"
// @Success 200 {object} object{message=string,user_id=string,username=string,email=string,cart_id=string} "User and cart created successfully"
// @Failure 400 {object} object{error=string} "Bad Request: Invalid input data"
// @Failure 409 {object} object{error=string} "Conflict: Username or email already exists"
// @Failure 500 {object} object{error=string} "Internal Server Error: Failed to create user or cart due to database/server error"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var userReq requests.CreateUser
	if err := c.ShouldBindJSON(&userReq); err != nil {
		log.Printf("Validation error for user creation request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting database transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate transaction for user creation."})
		return
	}

	newUserID := uuid.New().String()
	newCartID := uuid.New().String()

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		tx.Rollback()
		log.Printf("Error hashing password for user '%s': %v", userReq.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password for user creation."})
		return
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
		log.Printf("Error creating user '%s': %v", userReq.Username, userResult.Error)

		if errors.Is(userResult.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists. Please choose a different one."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user due to a database error."})
		}
		return
	}

	cartEntity := requests.Cart{
		Id:     newCartID,
		UserId: newUserID,
		Movies: []requests.CartMovieItem{},
	}

	cartResult := tx.Create(&cartEntity)
	if cartResult.Error != nil {
		tx.Rollback()
		log.Printf("Error creating cart for user '%s' (ID: %s): %v", userReq.Username, newUserID, cartResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user's cart. User creation rolled back."})
		return
	}

	tx.Commit()
	if tx.Error != nil {
		log.Printf("Error committing transaction for user '%s': %v", userReq.Username, tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize user and cart creation due to a commit error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User and cart created successfully!",
		"user_id":  userEntity.ID,
		"username": userEntity.Username,
		"email":    userEntity.Email,
		"cart_id":  cartEntity.Id,
	})
}

// LoginUser
// @Summary Authenticate user and get JWT token
// @Description Authenticates a user with email and password, returning a JWT token upon successful login.
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{email=string,password=string} true "User login credentials (email and password)"
// @Success 200 {object} object{message=string,token=string} "Login successful, returns JWT token"
// @Failure 400 {object} object{error=string} "Bad Request: Invalid input data"
// @Failure 401 {object} object{error=string} "Unauthorized: Incorrect password"
// @Failure 404 {object} object{error=string} "Not Found: User not found with provided email"
// @Failure 500 {object} object{error=string} "Internal Server Error: Failed to retrieve user or generate token"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var loginReq models.User
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Printf("Validation error for user creation request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := database.DB.Where("email = ?", loginReq.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User '%s' not found. Please ensure the user exists.", loginReq.Username)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User '%s' not found. Please ensure the user exists.", loginReq.Username)})
		} else {
			log.Printf("Database error retrieving user '%s': %v", loginReq.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user for login."})
		}
		return
	}

	passwordIsValid := utils.CheckPasswordHash(loginReq.Password, user.Password)
	if !passwordIsValid {
		log.Printf("Incorrect password for user '%s'", loginReq.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password."})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		log.Printf("Error generating token for user '%s': %v", loginReq.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token for login."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   token,
	})
}
