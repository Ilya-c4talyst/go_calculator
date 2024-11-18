package utils

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() (string, string) {
	runLocal := "1"
	port := "1"

	err := godotenv.Load()

	if err != nil {
		runLocal = "1"
		port = "8080"

	} else {
		runLocal = os.Getenv("runLocal")
		port = os.Getenv("port")
	}

	return runLocal, port
}
