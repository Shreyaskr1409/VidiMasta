package data

import "time"

// "gorm.io/gorm"
// "github.com/go-playground/validator/v10"

type User struct {
	Id           string
	Username     string
	Fullname     string
	Email        string
	Password     string
	Avatar       string
	RefreshToken string
	CreatedAt    *time.Time
}
