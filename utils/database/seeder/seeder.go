package seeder

import (
	"os"
	"tinder-clone/src/factory"
)

func Init(f *factory.Factory) {
	if os.Getenv("DB_SEEDER") == "false" {
		return
	}

	PremiumpackageSeeder(f.PremiumPackageRepository)
	UserSeeder(f.UserRepository)
}
