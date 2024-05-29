package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Base
	Name                    string `json:"name" gorm:"not null"`
	Email                   string `json:"email" gorm:"unique;not null"`
	Password                string `json:"password" gorm:"not null"`
	Age                     int    `json:"age" gorm:"not null"`
	IsExceedDailySwipeLimit *bool  `json:"is_exceed_daily_swipe_limit" gorm:"default:false;not null"`
	HasSwipeLimit           *bool  `json:"has_swipe_limit" gorm:"default:true;not null"`
	IsVerified              *bool  `json:"is_verified" gorm:"default:false;not null"`
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPasswordHash(inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
	return err == nil
}
