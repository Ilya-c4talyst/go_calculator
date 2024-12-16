package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHandler(t *testing.T) {
	// Создаем валидный/невалидный реквест
	data_s := []byte(`{"expression":"1+5"}`)
	suc := bytes.NewReader(data_s)
	data_f := []byte(`{"expression":"1/0"}`)
	fail := bytes.NewReader(data_f)
	// Пишем тест-кейсы
	testCases := []struct {
		name           string
		req            *http.Request
		expectedStatus int
		expected       float64
	}{
		{
			name:           "Success Case",
			req:            httptest.NewRequest(http.MethodPost, "/api/v1/calculate", suc),
			expectedStatus: http.StatusOK,
			expected:       6,
		},
		{
			name:           "Bad Request Case",
			req:            httptest.NewRequest(http.MethodPost, "/api/v1/calculate", fail),
			expectedStatus: http.StatusUnprocessableEntity,
			expected:       0,
		},
	}
	// Тестируем
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Запускаем хендлер
			w := httptest.NewRecorder()
			CalcHandler(w, tc.req)

			// Получаем результат
			res := w.Result()
			response := Response{}
			errorResponse := ErrorResponse{}

			// Получаем дату из респонса и десереализируем в структуру
			data, _ := io.ReadAll(res.Body)
			err := json.Unmarshal(data, &response)
			if err != nil {
				err = json.Unmarshal(data, &errorResponse)
			}
			defer res.Body.Close()

			// Проверки
			if err != nil {
				t.Error("Error while getting response")
			}

			if res.StatusCode != tc.expectedStatus {
				t.Errorf("Ошибка в статус коде: got %d, want %d", res.StatusCode, tc.expectedStatus)
			}

			if response.Result != tc.expected {
				t.Errorf("Ошибка в значении: got %f, want %f", response.Result, tc.expected)
			}
		})
	}
}
