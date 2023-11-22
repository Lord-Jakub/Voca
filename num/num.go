package num

import (
	"Voca/lib"
	"fmt"
)

const (
	INT      = 0
	OPERATOR = 1
)

func Evaluate(tokens map[int]string) (float64, error) {

	var stack []float64
	var operators []string

	precedence := map[string]float64{
		"*": 2,
		"/": 2,
		"+": 1,
		"-": 1,
	}

	for i := 0; i < len(tokens); i++ {
		tokenType := INT
		tokenValue := tokens[i]

		if i%2 == 1 {
			tokenType = OPERATOR
		}

		if tokenType == INT {

			num := lib.ParseFloat(tokenValue)

			stack = append(stack, num)
		} else if tokenType == OPERATOR {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[tokenValue] {
				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				b := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				a := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				result := 0.0

				switch op {
				case "+":
					result = a + b
				case "-":
					result = a - b
				case "*":
					result = a * b
				case "/":
					if b == 0 {
						return 0, fmt.Errorf("Dělení nulou")
					}
					result = a / b
				}
				stack = append(stack, result)
			}
			operators = append(operators, tokenValue)
		}
	}

	for len(operators) > 0 {
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]
		b := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		a := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result := 0.0

		switch op {
		case "+":
			result = a + b
		case "-":
			result = a - b
		case "*":
			result = a * b
		case "/":
			if b == 0 {
				return 0, fmt.Errorf("Dělení nulou")
			}
			result = a / b
		}
		stack = append(stack, result)
	}

	if len(stack) != 1 || len(operators) != 0 {
		return 0, fmt.Errorf("Neplatný výraz")
	}

	return stack[0], nil
}
