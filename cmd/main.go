package main

import (
	"github.com/Ilya-c4talyst/go_calculator/internal/application"
	"github.com/Ilya-c4talyst/go_calculator/internal/utils"
)

func main() {
	runLocal, port := utils.LoadEnv()
	app := application.NewApp(port)

	if runLocal == "1" {
		app.LocalRun()
	} else {
		app.RunServer()
	}
}
