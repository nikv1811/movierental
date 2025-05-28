package database

import (
	"fmt"
	"movierental/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	// dsn := "host=localhost user=nikhil.verma dbname=movie-rental-db port=5432 sslmode=disable"
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connected to database")
	}
}
