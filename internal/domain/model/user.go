package model

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Team     string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
