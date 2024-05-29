package dto

import "time"

type UserResponseDto struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Email                   string    `json:"email"`
	Age                     int       `json:"age"`
	IsExceedDailySwipeLimit bool      `json:"is_exceed_daily_swipe_limit"`
	HasSwipeLimit           bool      `json:"has_swipe_limit"`
	IsVerified              bool      `json:"is_verified"`
	CreatedAt               time.Time `json:"created_at"`
}

type UserFilterDto struct {
	Age            int
	MatchedUserIDs []string
}
