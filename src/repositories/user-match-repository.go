package repositories

import (
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"

	"gorm.io/gorm"
)

type UserMatchRepository struct {
	DB *gorm.DB
}

func NewUserMatchRepository(db *gorm.DB) *UserMatchRepository {
	return &UserMatchRepository{
		DB: db,
	}
}

func (u *UserMatchRepository) Create(c *abstraction.Context, userMatch *models.UserMatch) error {
	return u.DB.WithContext(c).Create(userMatch).Error
}

func (u *UserMatchRepository) FindMatchedUser(c *abstraction.Context, ID string) ([]string, error) {
	var matchedUserIDs []string
	err := u.DB.WithContext(c).Model(&models.UserMatch{}).Where("user_id = ? AND expire_time > NOW() - INTERVAL '24 HOURS'", ID).Pluck("match_user_id", &matchedUserIDs).Error
	return matchedUserIDs, err
}
