package routes

import (
	// _ "movierental/cmd/docs"
	"movierental/config" // Import config to access AppConfig.MovieAPI.BaseURL
	_ "movierental/docs"
	"movierental/pkg/controller"
	"movierental/pkg/middlewares"
	"movierental/pkg/movie/movieExternalApi" // Import the movieExternalApi package
	"movierental/pkg/services"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	userService := &services.UserService{}

	movieAPIClient := movieExternalApi.NewAPIClient(config.AppConfig.MovieAPI.BaseURL)

	movieService := services.NewMovieService(movieAPIClient)

	cartService := &services.CartService{}

	userController := &controller.UserController{UserService: userService}
	movieController := &controller.MovieController{MovieService: movieService}
	cartController := &controller.CartController{CartService: cartService}

	router.POST("/users", userController.CreateUser)
	router.POST("/login", userController.LoginUser)

	authenticatedGroup := router.Group("/")
	authenticatedGroup.Use(middlewares.Authenticate)
	{
		authenticatedGroup.GET("/listallmovies", movieController.ListAllMovies)
		authenticatedGroup.GET("/movie", movieController.MovieDetails)
		authenticatedGroup.GET("/cart", cartController.RetriveCart)
		authenticatedGroup.POST("/cart", cartController.AddToCart)
		authenticatedGroup.DELETE("/cart", cartController.RemoveFromCart)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
