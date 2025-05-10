package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Ilya-c4talyst/go_calculator/auth/cmd/auth_grpc"
	"github.com/Ilya-c4talyst/go_calculator/auth/database"
	"github.com/Ilya-c4talyst/go_calculator/auth/handlers"
	"github.com/Ilya-c4talyst/go_calculator/auth/middleware"
	"github.com/Ilya-c4talyst/go_calculator/auth/models"
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
		database.DB.AutoMigrate(&models.User{})
	}

	// Запускаем gRPC сервер в отдельной горутине
	go auth_grpc.StartGRPCServer(database.DB)

	// Создание роутера
	r := gin.Default()
	// Подключаем CORS middleware
	r.Use(middleware.CorsMiddleware())

	// Инициализация хендлеров
	authHandler := handlers.NewAuthHandler(database.DB)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Защищенны роуты
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			c.JSON(200, gin.H{"user_id": userID})
		})
	}

	// Запуск сервера
	r.Run(":8081")
}
