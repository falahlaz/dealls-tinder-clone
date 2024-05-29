package repositories

import (
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
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

func (u *UserRepository) Creates(c *abstraction.Context, users []*models.User) error {
	return u.DB.WithContext(c).Create(users).Error
}

func (u *UserRepository) FindByEmail(c *abstraction.Context, email string) (*models.User, error) {
	var user models.User
	err := u.DB.WithContext(c).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *UserRepository) FindByID(c *abstraction.Context, ID string) (*models.User, error) {
	var user models.User
	err := u.DB.WithContext(c).Where("id = ?", ID).First(&user).Error
	return &user, err
}

func (u *UserRepository) Find(c *abstraction.Context, filter *dto.UserFilterDto) ([]*models.User, error) {
	var users []*models.User
	err := u.DB.WithContext(c).Where("age BETWEEN ? AND ? AND id NOT IN ?", filter.Age-5, filter.Age, filter.MatchedUserIDs).Limit(20).Find(&users).Error
	return users, err
}

func (u *UserRepository) Update(c *abstraction.Context, user *models.User) error {
	return u.DB.WithContext(c).Where("id = ?", user.ID).Updates(user).Error
}
