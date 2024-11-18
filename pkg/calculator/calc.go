package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Инициализация ошибок
var ErrEmpty = errors.New("empty string")
var ErrPriority = errors.New("error with ()")
var ErrSymbols = errors.New("error with symbols")
var ErrZeroDivide = errors.New("zero divivsion")

func CheckPriorets(expression string) bool {
	// Валидация скобок
	stack := make([]string, 0)
	for _, v := range expression {

		if string(v) == "(" {
			stack = append(stack, string(v))
		} else if string(v) == ")" && len(stack) > 0 {
			stack = stack[:len(stack)-1]
		} else if string(v) == ")" {
			return false
		}
	}
	return len(stack) == 0
}

func CheckDigits(expression string) bool {
	// Валидация символов
	symbols := "+-/*()"
	// Строка начинается и заканчивается не на операнд
	if strings.HasPrefix(expression, "+") || strings.HasSuffix(expression, "+") ||
		strings.HasPrefix(expression, "-") || strings.HasSuffix(expression, "-") ||
		strings.HasPrefix(expression, "*") || strings.HasSuffix(expression, "*") ||
		strings.HasPrefix(expression, "/") || strings.HasSuffix(expression, "/") {
		return false
	}
	// Флаг, отвечающий за тип символа
	var flag string

	for _, v := range expression {
		if string(v) == "(" || string(v) == ")" {
			continue
		} else if unicode.IsDigit(rune(v)) && (flag == "" || flag == "symb") {
			flag = "dig"
		} else if !unicode.IsDigit(rune(v)) && (flag == "" || flag == "dig") && strings.Contains(symbols, string(v)) {
			flag = "symb"
		} else {
			return false
		}
	}

	return true
}

func Calc(expression string) (float64, error) {
	// Основная функция
	// Валидация входящей строки
	if len(expression) <= 2 {
		return 0, ErrEmpty
	}
	if strings.Contains(expression, "()") {
		return 0, ErrSymbols
	}
	if !CheckPriorets(expression) {
		return 0, ErrPriority
	}
	if !CheckDigits(expression) {
		return 0, ErrSymbols
	}
	if strings.Contains(expression, "/0") {
		return 0, ErrZeroDivide
	}

	noPrioritetsExp := expression
	noPrioritetsExp = strings.Replace(noPrioritetsExp, "(", " ", -1)
	noPrioritetsExp = strings.Replace(noPrioritetsExp, ")", " ", -1)

	splitExpression := strings.Split(noPrioritetsExp, " ")

	// Если не было знаков приоритета
	if len(splitExpression) == 1 {
		result, err := CalcMini(splitExpression[0])
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	// Иначе проходимся по каждому значению списка и делаем промежуточные вычисления
	for i, v := range splitExpression {

		if v == "" || len(v) == 1 ||
			strings.HasPrefix(v, "+") || strings.HasSuffix(v, "+") ||
			strings.HasSuffix(v, "/") || strings.HasPrefix(v, "/") ||
			strings.HasPrefix(v, "-") || strings.HasSuffix(v, "-") ||
			strings.HasPrefix(v, "*") || strings.HasSuffix(v, "*") {
			continue

		} else {
			r, err := CalcMini(v)
			if err != nil {
				return 0, err
			}
			splitExpression[i] = fmt.Sprintf("%1.0f", r)
		}
	}

	// Повторная итерация по приоритетам
	for i, v := range splitExpression {

		if (strings.HasPrefix(v, "+") || strings.HasPrefix(v, "-") ||
			strings.HasPrefix(v, "*") || strings.HasPrefix(v, "/")) && splitExpression[i-1] == "" {
			splitExpression[i-1] = string(v[0])
			splitExpression[i] = v[1:]

		} else if (strings.HasSuffix(v, "+") || strings.HasSuffix(v, "-") ||
			strings.HasSuffix(v, "*") || strings.HasSuffix(v, "/")) && splitExpression[i+1] == "" {
			splitExpression[i+1] = string(v[len(v)-1])
			splitExpression[i] = v[:len(v)-1]
		}
	}

	newExp := strings.Join(splitExpression, "")
	result, err := CalcMini(newExp)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func CalcMini(s string) (float64, error) {
	// Функция, которая проходит по строке, очищенной от приоритетов и вычисляет значения
	numbers := make([]float64, 0)
	symbols := make([]string, 0)

	for _, v := range s {
		// Наполняем списки с цифрами и операндами
		if unicode.IsDigit(rune(v)) {
			floa, _ := strconv.ParseFloat(string(v), 64)
			numbers = append(numbers, floa)
		} else {
			symbols = append(symbols, string(v))
		}
	}

	symbLen := len(symbols)
	res := 0.0

	// Пока есть операнды
	for symbLen > 0 {

		for i, v := range symbols {
			// Первым приоритетом пройдемся по умножениям и делениям
			if string(v) == "*" {
				// Умножение
				res = numbers[i] * numbers[i+1]
				numbers[i+1] = res
				if i < len(s)-1 {
					numbers = append(numbers[:i], numbers[i+1:]...)
					symbols = append(symbols[:i], symbols[i+1:]...)
				} else {
					symbols = symbols[:i]
					numbers = numbers[:i+1]
				}
				symbLen--
				break

			} else if string(v) == "/" {
				// Деление
				if numbers[i+1] == 0.0 {
					return 0, ErrZeroDivide
				}
				res = numbers[i] / numbers[i+1]
				numbers[i+1] = res
				if i < len(s)-1 {
					numbers = append(numbers[:i], numbers[i+1:]...)
					symbols = append(symbols[:i], symbols[i+1:]...)
				} else {
					symbols = symbols[:i]
					numbers = numbers[:i+1]
				}
				symbLen--
				break
			}
		}

		for i, v := range symbols {
			// Вторым приоритетом пройдемся по разностям и суммам
			if string(v) == "-" {
				// Разность
				res = numbers[i] - numbers[i+1]
				numbers[i+1] = res
				if i < len(s)-1 {
					numbers = append(numbers[:i], numbers[i+1:]...)
					symbols = append(symbols[:i], symbols[i+1:]...)
				} else {
					symbols = symbols[:i]
					numbers = numbers[:i+1]
				}
				symbLen--
				break

			} else if string(v) == "+" {
				// Сумма
				res = numbers[i] + numbers[i+1]
				numbers[i+1] = res
				if i < len(s)-1 {
					numbers = append(numbers[:i], numbers[i+1:]...)
					symbols = append(symbols[:i], symbols[i+1:]...)
				} else {
					symbols = symbols[:i]
					numbers = numbers[:i+1]
				}
				symbLen--
				break
			}
		}
	}

	return numbers[0], nil
}
