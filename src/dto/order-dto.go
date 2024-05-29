package dto

import "time"

type OrderRequestDto struct {
	PremiumPackageID string `json:"premium_package_id" validate:"required"`
}

type OrderResponseDto struct {
	ID               string                     `json:"id"`
	PremiumPackageID string                     `json:"premium_package_id"`
	PremiumPackage   *PremiumPackageResponseDto `json:"premium_package"`
	InvoiceNumber    string                     `json:"invoice_number"`
	TotalAmount      float64                    `json:"total_amount"`
	OrderTime        time.Time                  `json:"order_time"`
	PaymentTime      *time.Time                 `json:"payment_time"`
	ExpireTime       *time.Time                 `json:"expire_time"`
	Status           string                     `json:"status"`
}

type PremiumPackageResponseDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
