package utils

import (
	"movierental/pkg/models"
	"movierental/pkg/models/requests"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	db.AutoMigrate(&models.User{}, &requests.Cart{})
	return db
}

func ClearTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users;")
	db.Exec("DELETE FROM carts;")
}

func CreateTestUserAndCart(db *gorm.DB, userID string, email string) {
	user := models.User{ID: userID, Username: "testuser_" + userID, Email: email, Password: "hashedpassword"}
	db.Create(&user)
	cart := requests.Cart{Id: "cart-" + userID, UserId: userID, Movies: []requests.CartMovieItem{}}
	db.Create(&cart)
}
