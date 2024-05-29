package dto

type PotentialMatchRequestDto struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type MatchRequestDto struct {
	Direction     string `json:"direction" validate:"required,oneof=left right"`
	MatchedUserID string `json:"matched_user_id" validate:"required"`
}
