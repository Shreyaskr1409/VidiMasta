package data

import "time"

// "gorm.io/gorm"
// "github.com/go-playground/validator/v10"

type User struct {
	Id           string     `json:"_id"`
	Username     string     `json:"username"`
	Fullname     string     `json:"fullname"`
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	AvatarUrl    string     `json:"avatar_url"`
	RefreshToken string     `json:"refresh_token"`
	CreatedAt    *time.Time `json:"created_at"`
}
