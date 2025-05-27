package routes

import (
	// _ "movierental/cmd/docs"
	_ "movierental/docs"
	"movierental/pkg/controller"
	"movierental/pkg/middlewares"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	router.POST("/users", controller.CreateUser)
	router.POST("/login", controller.LoginUser)

	authenticatedGroup := router.Group("/")
	authenticatedGroup.Use(middlewares.Authenticate)
	{
		authenticatedGroup.GET("/listallmovies", controller.ListAllMovies)
		authenticatedGroup.GET("/movie", controller.MovieDetails)
		authenticatedGroup.GET("/cart", controller.RetriveCart)
		authenticatedGroup.POST("/cart", controller.AddToCart)
		authenticatedGroup.DELETE("/cart", controller.RemoveFromCart)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
