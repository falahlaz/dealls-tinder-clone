package repositories

import (
	"tinder-clone/src/abstraction"
	"tinder-clone/src/models"

	"gorm.io/gorm"
)

type PremiumPackageRepository struct {
	DB *gorm.DB
}

func NewPremiumPackageRepository(db *gorm.DB) *PremiumPackageRepository {
	return &PremiumPackageRepository{
		DB: db,
	}
}

func (p *PremiumPackageRepository) Create(c *abstraction.Context, premiumPackage *models.PremiumPackage) error {
	return p.DB.WithContext(c).Create(premiumPackage).Error
}

func (p *PremiumPackageRepository) FindByID(c *abstraction.Context, ID string) (*models.PremiumPackage, error) {
	var packageData models.PremiumPackage
	err := p.DB.WithContext(c).Where("id = ?", ID).First(&packageData).Error
	return &packageData, err
}
