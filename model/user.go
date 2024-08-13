package models

// UserScore represents the structure for user scores
type UserScore struct {
	Username string  `json:"username"`
	Points   float64 `json:"points"`
}
