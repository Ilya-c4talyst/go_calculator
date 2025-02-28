package application

// Сервер
type Application struct {
	Port string
}

// Агент
type Agent struct{}

// Реквест
type Request struct {
	Expression string `json:"expression"`
}
