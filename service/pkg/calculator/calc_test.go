package calculator

import "testing"

func TestValidateExpression(t *testing.T) {
	tests := []struct {
		expression string
		valid      bool
	}{
		{"2 + 2", true},
		{"2 + (3 * 4)", true},
		{"2 + (3 * 4))", false}, // Несогласованные скобки
		{"2 + a", false},        // Недопустимый символ
		{"", false},             // Пустое выражение
	}

	for _, test := range tests {
		err := ValidateExpression(test.expression)
		if (err == nil) != test.valid {
			t.Errorf("ValidateExpression(%s) = %v, ожидалось: %v", test.expression, err, test.valid)
		}
	}
}
