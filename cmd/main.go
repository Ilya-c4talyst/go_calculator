package main

import (
	"github.com/Ilya-c4talyst/go_calculator/internal/application"
)

func main() {

	// Запускаем сервер и агента
	app := application.NewApp("8080")
	agent := application.NewAgent()

	go agent.RunAgent()
	app.RunServer()

}
