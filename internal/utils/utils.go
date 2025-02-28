package utils

import "net/http"

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
