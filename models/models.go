package models

// Модель задачи
type Task struct {
	Id            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     rune    `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

// Модель выполненной задачи
type DoneTask struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

// Модель выражения
type Expression struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

// Респонс списка выражений
type ExpressionsResponse struct {
	Expressions []*Expression `json:"expressions"`
}

// Респонс выражения
type ExpressionResponse struct {
	Expression *Expression `json:"expression"`
}

type PostExpressionsResponse struct {
	Id int `json:"id"`
}

// Имитация БД (надо будет поменять на таблички)
var Tasks = []Task{}
var TasksDone = []DoneTask{}
var Expressions = []*Expression{}

// ID для выражений
var Id = 0
