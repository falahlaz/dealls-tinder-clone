package seeder

import (
	"context"
	"sync"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"
	"tinder-clone/src/repositories"

	"github.com/gofrs/uuid"
	"github.com/jaswdr/faker"
)

func UserSeeder(userRepository *repositories.UserRepository) {
	fakerInstance := faker.New()
	personFaker := fakerInstance.Person()
	UUIDFaker := fakerInstance.UUID()
	cc := &abstraction.Context{
		Context: context.Background(),
	}

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	userFakers := []*models.User{}
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			userModel := &models.User{
				Base: models.Base{
					ID: uuid.Must(uuid.FromString(UUIDFaker.V4())),
				},
				Name:     personFaker.Name(),
				Email:    personFaker.Contact().Email,
				Password: personFaker.Faker.Internet().Password(),
				Age:      fakerInstance.IntBetween(18, 60),
			}
			userModel.HashPassword()

			mu.Lock()
			userFakers = append(userFakers, userModel)
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	if err := userRepository.Creates(cc, userFakers); err != nil {
		panic(err)
	}
}
