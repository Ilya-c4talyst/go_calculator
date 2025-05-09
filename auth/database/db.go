package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB переменная
var DB *gorm.DB

// Инициализация БД
func InitDatabase(POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_NAME, POSTGRES_PORT string) error {

	// Подключение к БД
	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_NAME, POSTGRES_PORT)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
