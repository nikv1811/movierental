package requests

type Cart struct {
	Id     string          `gorm:"primaryKey"`
	UserId string          `json:"user_id"`
	Movies []CartMovieItem `json:"movies" gorm:"type:json;serializer:json"`
}
