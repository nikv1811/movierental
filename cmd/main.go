package main

import (
	"movierental/pkg/database"
	"movierental/pkg/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	database.ConnectToDb()
}

// Swagger API Documentation Route
// @title Movie Rental API
// @version 1.0
// @description This is a Movie Rental API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8080")
}
