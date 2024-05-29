package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
