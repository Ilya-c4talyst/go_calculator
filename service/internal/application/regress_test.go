package application

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/Ilya-c4talyst/go_calculator/service/models"
// )

// // Тест для POST /api/v1/calculate
// func TestCalculateHandler(t *testing.T) {
// 	// Создаем тестовый сервер
// 	server := httptest.NewServer(http.HandlerFunc(CalcHandler))
// 	defer server.Close()

// 	// Тестовые данные
// 	requestBody := map[string]string{
// 		"expression": "2+2",
// 	}
// 	jsonData, err := json.Marshal(requestBody)
// 	if err != nil {
// 		t.Fatalf("Ошибка при кодировании JSON: %v", err)
// 	}

// 	// Выполняем POST-запрос
// 	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatalf("Ошибка при выполнении запроса: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Проверяем статус код
// 	if resp.StatusCode != http.StatusCreated {
// 		t.Errorf("Ожидался статус код 201, получен %d", resp.StatusCode)
// 	}

// 	// Проверяем тело ответа
// 	var response map[string]int
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		t.Fatalf("Ошибка при декодировании JSON: %v", err)
// 	}

// 	if _, ok := response["id"]; !ok {
// 		t.Error("Ожидалось поле 'id' в ответе")
// 	}
// }

// // Тест для GET /api/v1/expressions
// func TestExpressionsHandler(t *testing.T) {
// 	// Создаем тестовый сервер
// 	server := httptest.NewServer(http.HandlerFunc(ExpressionsHandler))
// 	defer server.Close()

// 	// Выполняем GET-запрос
// 	resp, err := http.Get(server.URL + "/api/v1/expressions")
// 	if err != nil {
// 		t.Fatalf("Ошибка при выполнении запроса: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Проверяем статус код
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Ожидался статус код 200, получен %d", resp.StatusCode)
// 	}

// 	// Проверяем тело ответа
// 	var response models.ExpressionsResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		t.Fatalf("Ошибка при декодировании JSON: %v", err)
// 	}

// 	if len(response.Expressions) == 0 {
// 		t.Error("Ожидался непустой список выражений")
// 	}
// }

// // Тест для GET /api/v1/expressions/:id
// func TestExpressionByIDHandler(t *testing.T) {
// 	// Создаем тестовый сервер
// 	server := httptest.NewServer(http.HandlerFunc(ExpressionHandler))
// 	defer server.Close()

// 	// Выполняем GET-запрос
// 	time.Sleep(time.Second * 5)
// 	resp, err := http.Get(server.URL + "/api/v1/expression/0")
// 	if err != nil {
// 		t.Fatalf("Ошибка при выполнении запроса: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Проверяем статус код
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Ожидался статус код 200, получен %d", resp.StatusCode)
// 	}

// 	// Проверяем тело ответа
// 	var response struct {
// 		Expression models.Expression `json:"expression"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		t.Fatalf("Ошибка при декодировании JSON: %v", err)
// 	}

// 	if response.Expression.Id != 0 {
// 		t.Errorf("Ожидался ID 0, получен %d", response.Expression.Id)
// 	}
// }
