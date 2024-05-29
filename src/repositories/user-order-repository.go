package repositories

import (
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"

	"gorm.io/gorm"
)

type UserOrderRepository struct {
	DB *gorm.DB
}

func NewUserOrderRepository(db *gorm.DB) *UserOrderRepository {
	return &UserOrderRepository{
		DB: db,
	}
}

func (u *UserOrderRepository) CountOrders(c *abstraction.Context) (int64, error) {
	var count int64
	err := u.DB.WithContext(c).Model(&models.UserOrder{}).Count(&count).Error
	return count, err
}
func (u *UserOrderRepository) Create(c *abstraction.Context, userOrder *models.UserOrder) error {
	return u.DB.WithContext(c).Create(userOrder).Error
}

func (u *UserOrderRepository) FindByID(c *abstraction.Context, ID string) (*models.UserOrder, error) {
	var userOrder models.UserOrder
	err := u.DB.WithContext(c).Preload("Package").Where("id = ?", ID).First(&userOrder).Error
	return &userOrder, err
}

func (u *UserOrderRepository) Update(c *abstraction.Context, userOrder *models.UserOrder) error {
	return u.DB.WithContext(c).Where("id = ?", userOrder.ID.String()).Updates(userOrder).Error
}
