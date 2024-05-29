package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserMatch struct {
	Base
	UserID      uuid.UUID `json:"user_id" gorm:"not null"`
	MatchUserID uuid.UUID `json:"match_user_id" gorm:"not null"`
	IsMatch     bool      `json:"is_match" gorm:"not null"`
	ExpireTime  time.Time `json:"expire_time" gorm:"default:null"`

	User  *User `json:"user" gorm:"foreignKey:UserID"`
	Match *User `json:"match" gorm:"foreignKey:MatchUserID"`
}
