package application

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ilya-c4talyst/go_calculator/service/database"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/auth_client"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/utils"
	"github.com/Ilya-c4talyst/go_calculator/service/models"
	"github.com/Ilya-c4talyst/go_calculator/service/pkg/calculator"
)

// Создание приложения
func NewApp(port string, authCli *auth_client.Client) *Application {
	return &Application{
		Port:    port,
		AuthCli: authCli,
	}
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

	userID := r.Context().Value("userID").(uint32)

	// Создаем задачи для выражений и заносим в имитацию БД
	expression := models.Expression{}
	expression.Value = request.Expression
	expression.Status = "queued"
	expression.Result = "-"
	expression.User_id = userID
	err = database.AddExpression(&expression)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	go func() {
		solved := calculator.EvaluateExpression(request.Expression, &expression)
		database.UpdateExpression(solved)
	}()
	log.Println("Calc started")

	// Отдаем пользователю информацию о созданной задаче
	response := models.PostExpressionsResponse{
		Id: expression.Id,
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

	userID := r.Context().Value("userID").(uint32)
	expressions, err := database.GetExpressions(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response := models.ExpressionsResponse{
		Expressions: expressions,
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

	expression, err := database.GetExpressionbyID(r.Context().Value("userID").(uint32), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

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
	protected := utils.AuthMiddleware(a.AuthCli)

	mux.Handle("/api/v1/calculate", protected(http.HandlerFunc(CalcHandler)))
	mux.Handle("/api/v1/expression/", protected(http.HandlerFunc(ExpressionHandler)))
	mux.Handle("/api/v1/expressions", protected(http.HandlerFunc(ExpressionsHandler)))
	mux.HandleFunc("/api/v1/internal/task", InternalTaskHandler) // без аутентификации

	handler := utils.EnableCORS(mux)

	log.Println("Server started")
	err := http.ListenAndServe(":"+a.Port, handler)

	if err != nil {
		log.Println("Error while running server")
	}

	return err
}
