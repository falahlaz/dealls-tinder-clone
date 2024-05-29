package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type UserOrder struct {
	Base
	PremiumPackageID uuid.UUID  `json:"premium_package_id" gorm:"not null"`
	UserID           uuid.UUID  `json:"user_id" gorm:"not null"`
	InvoiceNumber    string     `json:"invoice_number" gorm:"unique;not null"`
	TotalAmount      float64    `json:"total_amount" gorm:"not null"`
	OrderTime        time.Time  `json:"order_time" gorm:"not null"`
	PaymentTime      *time.Time `json:"payment_time" gorm:"default:null"`
	ExpireTime       *time.Time `json:"expire_time" gorm:"default:null"`
	Status           string     `json:"status" gorm:"not null"`

	Package *PremiumPackage `json:"package" gorm:"foreignKey:PremiumPackageID"`
	User    *User           `json:"user" gorm:"foreignKey:UserID"`
}

func (u *UserOrder) GenerateInvoiceNumber(count int64) {
	u.InvoiceNumber = fmt.Sprintf("ORDER-INV-%04d", count+1)
}
