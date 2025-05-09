package utils

import (
	"errors"
	"os"
	"time"

	models "github.com/Ilya-c4talyst/go_calculator/auth/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Конфиг JWT
type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

// Конструктор конфига
func GetJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey: os.Getenv("SECRETKEY"),
		ExpiresIn: 24 * time.Hour,
	}
}

// Генерация токена JWT
func GenerateToken(userID uint) (string, error) {
	cfg := GetJWTConfig()

	claims := models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// Парсинг токена
func ParseToken(tokenString string) (*models.Claims, error) {

	cfg := GetJWTConfig()

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка хеша
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
