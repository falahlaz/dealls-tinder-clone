package repositories

import (
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) Create(c *abstraction.Context, user *models.User) error {
	return u.DB.WithContext(c).Create(user).Error
}

func (u *UserRepository) FindByEmail(c *abstraction.Context, email string) (*models.User, error) {
	var user models.User
	err := u.DB.WithContext(c).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *UserRepository) Update(c *abstraction.Context, user *models.User) error {
	return u.DB.WithContext(c).Where("id = ?", user.ID).Updates(user).Error
}
