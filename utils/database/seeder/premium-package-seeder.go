package seeder

import (
	"context"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"
	"tinder-clone/src/repositories"

	"github.com/gofrs/uuid"
)

func PremiumpackageSeeder(repository *repositories.PremiumPackageRepository) {
	packages := []*models.PremiumPackage{
		{
			Base: models.Base{
				ID: uuid.Must(uuid.FromString("50f484f4-2ff8-4dbe-b5c6-105cd28186de")),
			},
			Name:        "No Swipe Quota",
			Amount:      100000,
			Description: "No Swipe Quota",
		},
		{
			Base: models.Base{
				ID: uuid.Must(uuid.FromString("5d56f9de-bded-4df9-99f7-4209a6b58dbe")),
			},
			Name:        "Verified Label",
			Amount:      100000,
			Description: "Verified Label",
		},
	}

	cc := &abstraction.Context{
		Context: context.Background(),
	}

	for _, packageData := range packages {
		if err := repository.Create(cc, packageData); err != nil {
			panic(err)
		}
	}
}
