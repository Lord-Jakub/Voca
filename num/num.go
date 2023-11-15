package num

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	INT      = 0
	OPERATOR = 1
)

func Evaluate(tokens map[int]string) (int, error) {
	n := 0
	for n < len(tokens) {

		if strings.Contains(tokens[n], ":") {
			tokens[n] = strings.SplitN(tokens[n], ":", 2)[1]
		}
		n++
	}

	var stack []int
	var operators []string

	precedence := map[string]int{
		"MULT":  2,
		"DIV":   2,
		"PLUS":  1,
		"MINUS": 1,
	}

	for i := 0; i < len(tokens); i++ {
		tokenType := INT
		tokenValue := tokens[i]

		if i%2 == 1 {
			tokenType = OPERATOR
		}

		if tokenType == INT {
			num, err := strconv.Atoi(tokenValue)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if tokenType == OPERATOR {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[tokenValue] {
				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				b := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				a := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				result := 0

				switch op {
				case "PLUS":
					result = a + b
				case "MINUS":
					result = a - b
				case "MULT":
					result = a * b
				case "DIV":
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
		result := 0

		switch op {
		case "PLUS":
			result = a + b
		case "MINUS":
			result = a - b
		case "MULT":
			result = a * b
		case "DIV":
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
