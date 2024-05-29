package database

import (
	"os"
	"tinder-clone/src/models"

	"gorm.io/gorm"
)

var (
	DBConnection *gorm.DB
)

func Init() {
	dbConfig := &databaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SslMode:  os.Getenv("DB_SSL_MODE"),
		Tz:       os.Getenv("DB_TZ"),
		LogLevel: os.Getenv("DB_LOG_LEVEL"),
	}

	var err error
	DBConnection, err = dbConfig.InitDB()
	if err != nil {
		panic(err)
	}

	if os.Getenv("DB_MIGRATE") == "true" {
		DBConnection.AutoMigrate(
			&models.PremiumPackage{},
			&models.User{},
			&models.UserOrder{},
			&models.UserMatch{},
		)
	}
}
