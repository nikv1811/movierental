package controller

import (
	"log"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"movierental/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserServiceInterface
}

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
func (uc *UserController) CreateUser(c *gin.Context) {
	var userReq requests.CreateUser
	if err := c.ShouldBindJSON(&userReq); err != nil {
		log.Printf("Validation error for user creation request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.UserService.CreateUser(userReq)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if err.Error() == "username or email already exists. Please choose a different one" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response)
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
func (uc *UserController) LoginUser(c *gin.Context) {
	var loginReq models.User
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Printf("Validation error for user login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.UserService.LoginUser(loginReq)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		if err.Error() == "incorrect password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if err.Error() == "user '' not found. Please ensure the user exists" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, response)
}
