package main

import (
	"log"
	"os"

	"github.com/Ilya-c4talyst/go_calculator/service/database"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/application"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/auth_client"
	"github.com/Ilya-c4talyst/go_calculator/service/models"
	"github.com/joho/godotenv"
)

func main() {
	// Получение данных окружения
	godotenv.Load()
	var POSTGRES_USER = os.Getenv("POSTGRES_USER")
	var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	var POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	var POSTGRES_NAME = os.Getenv("POSTGRES_NAME")
	var POSTGRES_HOST = os.Getenv("POSTGRES_HOST")

	// Инициализация базы данных
	errDatabase := database.InitDatabase(POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_NAME, POSTGRES_PORT)
	if errDatabase != nil {
		log.Fatal("Ошибка подключения к базе данных: ", errDatabase)
	} else {
		log.Println("Успешное подключение к базе данных")
		// Автоматическое создание таблицы на основе модели Note, если она не существует
		database.DB.AutoMigrate(&models.Expression{})
	}

	// Инициализация gRPC клиента
	authCli, err := auth_client.New("auth_service:50051")
	if err != nil {
		log.Fatal("Failed to create auth client: ", err)
	}
	defer authCli.Close()

	app := application.NewApp("8080", authCli)
	agent := application.NewAgent(authCli)

	go agent.RunAgent()

	if err := app.RunServer(); err != nil {
		log.Fatal("Server error: ", err)
	}
}
