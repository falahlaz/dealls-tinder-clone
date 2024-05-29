package factory

import (
	"fmt"
	"os"
	"tinder-clone/src/repositories"
	"tinder-clone/utils/database"

	"github.com/redis/go-redis/v9"
)

type Factory struct {
	PremiumPackageRepository *repositories.PremiumPackageRepository
	UserMatchRepository      *repositories.UserMatchRepository
	UserOrderRepository      *repositories.UserOrderRepository
	UserRepository           *repositories.UserRepository

	RedisClient *redis.Client
}

func NewFactory() *Factory {
	return &Factory{
		PremiumPackageRepository: repositories.NewPremiumPackageRepository(database.DBConnection),
		UserMatchRepository:      repositories.NewUserMatchRepository(database.DBConnection),
		UserOrderRepository:      repositories.NewUserOrderRepository(database.DBConnection),
		UserRepository:           repositories.NewUserRepository(database.DBConnection),

		RedisClient: InitRedis(),
	}
}

func InitRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
}
