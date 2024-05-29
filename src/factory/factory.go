package factory

import (
	"tinder-clone/src/repositories"
	"tinder-clone/utils/database"
)

type Factory struct {
	PremiumPackageRepository *repositories.PremiumPackageRepository
	UserOrderRepository      *repositories.UserOrderRepository
	UserRepository           *repositories.UserRepository
}

func NewFactory() *Factory {
	return &Factory{
		PremiumPackageRepository: repositories.NewPremiumPackageRepository(database.DBConnection),
		UserOrderRepository:      repositories.NewUserOrderRepository(database.DBConnection),
		UserRepository:           repositories.NewUserRepository(database.DBConnection),
	}
}
