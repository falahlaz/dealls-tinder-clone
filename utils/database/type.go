package database

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db interface {
	InitDB() (*gorm.DB, error)
}

type databaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SslMode  string
	Tz       string
	LogLevel string
}

func (dc *databaseConfig) InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dc.Host, dc.Username, dc.Password, dc.Name, dc.Port, dc.SslMode, dc.Tz)

	// envLogLevelRead is log for all database
	envLogLevelWrite := strings.ToLower(dc.LogLevel)
	logLevelWrite := setLogLevel(envLogLevelWrite)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevelWrite),
	}) // default connection with write db
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setLogLevel(param string) logger.LogLevel {
	switch param {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "silent":
		return logger.Silent
	}

	return logger.Warn // warn is default gorm log level
}
