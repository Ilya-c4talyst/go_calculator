package utils

import (
	"context"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/auth_client"
	"net/http"
)

// Семафор
type Semaphore struct {
	n  int
	ch chan int
}

// NewSemaphore создает новый семафор указанной вместимости.
func NewSemaphore(n int) *Semaphore {
	c := make(chan int, n)
	return &Semaphore{n, c}
}

// Acquire занимает место в семафоре, если есть свободное.
func (s *Semaphore) Acquire() {
	s.ch <- 1
}

// Release освобождает место в семафоре
func (s *Semaphore) Release() {
	<-s.ch
}

// Middleware для обработки CORS
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешаем определенные методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Разрешаем определенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос дальше
		next.ServeHTTP(w, r)
	})
}

// Middleware для проверки прав
func AuthMiddleware(authCli *auth_client.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен аутентификации
			token := r.Header.Get("Authorization")

			// Валидируем его и получаем данные о пользователе
			if token == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			userID, err := authCli.ValidateToken(token)

			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Создаем новый контекст с userID
			ctx := context.WithValue(r.Context(), "userID", userID)

			// Создаем новый запрос с обновленным контекстом
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
