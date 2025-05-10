package application

import "github.com/Ilya-c4talyst/go_calculator/service/internal/auth_client"

// Сервер
type Application struct {
	Port    string
	AuthCli *auth_client.Client
}

// Агент
type Agent struct {
	AuthCli *auth_client.Client
}

// Реквест
type Request struct {
	Expression string `json:"expression"`
}
