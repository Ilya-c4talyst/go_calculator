package calculator

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/Ilya-c4talyst/go_calculator/models"
	"github.com/joho/godotenv"
)

// ID таски (TODO: перенести в БД)
var id = 1

// Валидация входящего выражения
func ValidateExpression(expression string) error {

	// Удаляем все пробелы
	expression = strings.ReplaceAll(expression, " ", "")

	if expression == "" {
		return fmt.Errorf("пустое выражение")
	}

	// Проверка на допустимые символы
	for _, char := range expression {
		if !unicode.IsDigit(char) && char != '.' && char != '+' && char != '-' && char != '*' && char != '/' && char != '(' && char != ')' {
			return fmt.Errorf("недопустимый символ: %c", char)
		}
	}

	// Проверка сбалансированности скобок
	balance := 0
	for _, char := range expression {
		if char == '(' {
			balance++
		} else if char == ')' {
			balance--
			if balance < 0 {
				return fmt.Errorf("несогласованные скобки")
			}
		}
	}
	if balance != 0 {
		return fmt.Errorf("несогласованные скобки")
	}

	return nil
}

// Вычисление итогового значения
func EvaluateExpression(expression string, current *models.Expression) {

	// Получение переменных среды
	godotenv.Load()
	var TIME_ADDITION_MS, _ = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	var TIME_SUBTRACTION_MS, _ = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	var TIME_MULTIPLICATIONS_MS, _ = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	var TIME_DIVISIONS_MS, _ = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))

	// Удаляем все пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	// Стек для чисел
	var numbers []float64
	// Стек для операторов
	var operators []rune

	// Функция для выполнения операции
	applyOperation := func(op rune) error {
		if len(numbers) < 2 {
			return fmt.Errorf("недостаточно операндов для операции %c", op)
		}
		b := numbers[len(numbers)-1]
		a := numbers[len(numbers)-2]
		numbers = numbers[:len(numbers)-2]

		// Формируем задачу для агента
		var result float64
		var task models.Task

		task.Id = id
		id++
		task.Arg1 = a
		task.Arg2 = b
		task.Operation = op

		switch op {
		case '+':
			task.OperationTime = TIME_ADDITION_MS
		case '-':
			task.OperationTime = TIME_SUBTRACTION_MS
		case '*':
			task.OperationTime = TIME_MULTIPLICATIONS_MS
		case '/':
			task.OperationTime = TIME_DIVISIONS_MS
		}

		models.Tasks = append(models.Tasks, task)
		log.Println("Task created")

		// Ждем ответ от агента, проверяя выполненные задачи

	Loop:
		for {
			for _, t := range models.TasksDone {
				if t.Id == task.Id {
					result = t.Result
					models.TasksDone = append(models.TasksDone[:len(models.TasksDone)-1], models.TasksDone[len(models.TasksDone):]...)
					break Loop
				}
			}
		}

		log.Println("Loop closed, tasks done")

		numbers = append(numbers, result)
		return nil
	}

	// Обходим каждый символ в выражении
	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' {
			// Если символ - цифра или точка, собираем число
			j := i
			for j < len(expression) && (unicode.IsDigit(rune(expression[j])) || expression[j] == '.') {
				j++
			}
			num, err := strconv.ParseFloat(expression[i:j], 64)

			if err != nil {
				current.Result = fmt.Sprintf("ошибка при парсинге числа: %v", err)
				current.Status = "broken"
				return
			}

			numbers = append(numbers, num)
			i = j - 1

		} else if char == '(' {
			// Если символ - открывающая скобка, добавляем в стек операторов
			operators = append(operators, char)

		} else if char == ')' {
			// Если символ - закрывающая скобка, выполняем операции до открывающей скобки
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if err := applyOperation(operators[len(operators)-1]); err != nil {
					current.Result = err.Error()
					current.Status = "broken"
					return
				}
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				current.Result = "несогласованные скобки"
				current.Status = "broken"
				return
			}
			// Убираем открывающую скобку
			operators = operators[:len(operators)-1]

		} else if char == '+' || char == '-' || char == '*' || char == '/' {

			// Если символ - оператор, выполняем операции с более высоким приоритетом
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(char) {
				if err := applyOperation(operators[len(operators)-1]); err != nil {
					current.Result = err.Error()
					current.Status = "broken"
					return
				}
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)

		} else {
			current.Result = fmt.Sprintf("неизвестный символ: %c", char)
			current.Status = "broken"
			return
		}
	}

	// Выполняем оставшиеся операции
	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			current.Result = "несогласованные скобки"
			current.Status = "broken"
			return
		}
		if err := applyOperation(operators[len(operators)-1]); err != nil {
			current.Result = err.Error()
			current.Status = "broken"
			return
		}
		operators = operators[:len(operators)-1]
	}

	if len(numbers) != 1 {
		current.Result = "ошибка в выражении"
		current.Status = "broken"
		return
	}

	current.Result = strconv.FormatFloat(numbers[0], 'f', 2, 64)
	current.Status = "solved"
	log.Println("Expression evaluated successfully")
}

// Функция для определения приоритета оператора
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}
