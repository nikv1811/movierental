package main

import (
	"movierental/pkg/database"
	"movierental/pkg/models"
	"movierental/pkg/models/requests"
)

func init() {
	database.ConnectToDb()
}

func main() {
	database.DB.AutoMigrate(&models.User{}, &requests.Cart{})
}
