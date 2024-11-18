package calculator

import (
	"testing"
)

func TestCalc(t *testing.T) {

	cases := []struct {
		name  string
		input string
		want  float64
		er    error
	}{

		{
			name:  "1",
			input: "5*(6-6*1)",
			want:  0,
			er:    nil,
		},

		{
			name:  "2",
			input: "8/2*(6+1)",
			want:  28,
			er:    nil,
		},
		{
			name:  "3",
			input: "1/2*(6+1)",
			want:  3.5,
			er:    nil,
		},
		{
			name:  "4",
			input: "+7-9+2+1-2+8+1-2",
			want:  0.0,
			er:    ErrSymbols,
		},
		{
			name:  "5",
			input: "1+9",
			want:  10,
			er:    nil,
		},
		{
			name:  "6",
			input: "5*3/(5-5)",
			want:  0.0,
			er:    ErrZeroDivide,
		},
		{
			name:  "7",
			input: "5*3/((5-5)",
			want:  0.0,
			er:    ErrPriority,
		},
		{
			name:  "8",
			input: "",
			want:  0.0,
			er:    ErrEmpty,
		},
	}
	// перебор всех тестов
	for _, tc := range cases {

		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, er := Calc(tc.input)
			if got != tc.want && er != tc.er {
				t.Errorf("Error in test %s", tc.name)
			}
		})
	}
}
