package main

import (
	"movierental/pkg/controller"
	"movierental/pkg/database"
	"movierental/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	database.ConnectToDb()
}

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	router.POST("/users", controller.CreateUser)
	router.POST("/login", controller.LoginUser)

	authenticatedGroup := router.Group("/")
	authenticatedGroup.Use(middlewares.Authenticate)
	authenticatedGroup.GET("/listallmovies", controller.ListAllMovies)
	authenticatedGroup.GET("/movie", controller.MovieDetails)
	authenticatedGroup.GET("/cart/:user_id", controller.RetriveCart)
	authenticatedGroup.POST("/cart/:user_id", controller.AddToCart)

	router.Run(":8080")
}
