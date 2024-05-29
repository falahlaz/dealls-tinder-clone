package models

type PremiumPackage struct {
	Base
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
}
