package application

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ilya-c4talyst/go_calculator/internal/utils"
	"github.com/Ilya-c4talyst/go_calculator/models"
	"github.com/Ilya-c4talyst/go_calculator/pkg/calculator"
)

// Создание приложения
func NewApp(port string) *Application {
	return &Application{Port: port}
}

// Главный обработчик для выражений
func CalcHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Получение выражения из реквеста
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		log.Println("Error while using calculator:", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	// Валидация выражения
	if err := calculator.ValidateExpression(request.Expression); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	log.Println("Validation passed")

	// Создаем задачи для выражений и заносим в имитацию БД
	Expression := models.Expression{Id: models.Id}
	Expression.Status = "queued"
	Expression.Result = "-"
	models.Expressions = append(models.Expressions, &Expression)

	go func() {
		models.Id++
		calculator.EvaluateExpression(request.Expression, &Expression)
	}()
	log.Println("Calc started")

	// Отдаем пользователю информацию о созданной задаче
	response := models.PostExpressionsResponse{
		Id: Expression.Id,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		log.Println("Error while marshal JSON", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Хэндлер для получения всех выражений
func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := models.ExpressionsResponse{
		Expressions: models.Expressions,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error while marshal JSON", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Хэндлер для получения одного выражений по ключу
func ExpressionHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/api/v1/expression/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 0 || id >= len(models.Expressions) {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	expression := models.Expressions[id]

	response := models.ExpressionResponse{
		Expression: expression,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error while marshal JSON", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Хэндлер для работы с задачами (используется агентом)
func InternalTaskHandler(w http.ResponseWriter, r *http.Request) {

	// Гет запрос на получение задачи
	if r.Method == "GET" {

		w.Header().Set("Content-Type", "application/json")

		if len(models.Tasks) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(models.Tasks[0])
		models.Tasks = models.Tasks[1:]

		if err != nil {
			log.Println("Error while marshal JSON", err)
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		// Пост запрос
	} else if r.Method == "POST" {

		var result models.DoneTask

		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		models.TasksDone = append(models.TasksDone, result)
		w.WriteHeader(http.StatusCreated)
	}
}

// Запуск сервера
func (a *Application) RunServer() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", CalcHandler)
	mux.HandleFunc("/api/v1/expression/", ExpressionHandler)
	mux.HandleFunc("/api/v1/expressions", ExpressionsHandler)
	mux.HandleFunc("/api/v1/internal/task", InternalTaskHandler)

	handler := utils.EnableCORS(mux)

	log.Println("Server started")
	err := http.ListenAndServe(":"+a.Port, handler)

	if err != nil {
		log.Println("Error while running server")
	}

	return err
}
