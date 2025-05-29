package services

import (
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
	"movierental/pkg/movie/movieExternalApi"
)

type UserServiceInterface interface {
	CreateUser(userReq requests.CreateUser) (map[string]interface{}, error)
	LoginUser(loginReq models.User) (map[string]interface{}, error)
}

type MovieServiceInterface interface {
	ListAllMovies(queryParams map[string]string) ([]movieExternalApi.Movie, error)
	GetMovieDetails(movieId string) (movieExternalApi.Movie, error)
}

type CartServiceInterface interface {
	RetrieveCart(userId interface{}) (requests.Cart, error)
	AddToCart(userId interface{}, movieItem requests.CartMovieItem) (map[string]interface{}, error)
	RemoveFromCart(userId interface{}, movieID int) (map[string]interface{}, error)
}
