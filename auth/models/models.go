package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Модель пользователя
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Креды
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Ответ от аутентификации (токен)
type AuthResponse struct {
	Token string `json:"token"`
}
