package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Ilya-c4talyst/go_calculator/service/internal/auth_client"
	"github.com/Ilya-c4talyst/go_calculator/service/internal/utils"
	"github.com/Ilya-c4talyst/go_calculator/service/models"
	"github.com/joho/godotenv"
)

// Создание агента
func NewAgent(authCli *auth_client.Client) *Agent {
	return &Agent{AuthCli: authCli}
}

func (a *Agent) RunAgent() {

	log.Println("Agent Started")

	// Подгрузка переменных среды
	godotenv.Load()
	var TIME_ADDITION_MS, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	var TIME_SUBTRACTION_MS, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	var TIME_MULTIPLICATIONS_MS, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	var TIME_DIVISIONS_MS, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	var COMPUTING_POWER, _ = strconv.Atoi(os.Getenv("COMPUTING_POWER"))

	// Создание мьютекса и семафора
	var mutex sync.Mutex
	sema := utils.NewSemaphore(COMPUTING_POWER)

	// Запускаем агента, который будет постоянно опрашивать наш сервер
	for {
		sema.Acquire()

		// Основная логика агента
		go func() {
			mutex.Lock()
			defer sema.Release()
			defer mutex.Unlock()

			var task models.Task
			resp, err := http.Get("http://127.0.0.1:8080/api/v1/internal/task")

			if err != nil {
				log.Println(err, "Get task error")
				return
			}

			err = json.NewDecoder(resp.Body).Decode(&task)

			if err != nil {
				return
			}

			// Получили таску - обрабатываем
			op := string(task.Operation)
			a, b := task.Arg1, task.Arg2
			result := 0.0

			switch op {
			case "+":
				time.Sleep(time.Duration(TIME_ADDITION_MS))
				result = a + b
			case "-":
				result = a - b
				time.Sleep(time.Duration(TIME_SUBTRACTION_MS))
			case "*":
				time.Sleep(time.Duration(TIME_MULTIPLICATIONS_MS))
				result = a * b
			case "/":
				time.Sleep(time.Duration(TIME_DIVISIONS_MS))
				result = a / b
			}

			taskDone := models.DoneTask{Id: task.Id, Result: result}
			jsonBytes, _ := json.Marshal(taskDone)

			// Отправляем результат
			_, err = http.Post("http://127.0.0.1:8080/api/v1/internal/task", "application/json", bytes.NewReader(jsonBytes))
			if err != nil {
				log.Println(err, "Post task")
				return
			}
		}()
	}
}
