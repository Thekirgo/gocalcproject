package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := Tokenize(expression)
	rpn, err := infixToRPN(tokens)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpn)
}

func Tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder

	for _, char := range expression {
		switch char {
		case ' ':
			continue
		case '+', '-', '*', '/', '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(char))
		default:
			number.WriteRune(char)
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens
}

func infixToRPN(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		switch token {
		case "+":
			for len(stack) > 0 && (stack[len(stack)-1] == "+" || stack[len(stack)-1] == "-") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "-":
			for len(stack) > 0 && (stack[len(stack)-1] == "+" || stack[len(stack)-1] == "-") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "*":
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "/":
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		default:
			output = append(output, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token {
		case "+":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a+b)
		case "-":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a-b)
		case "*":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a*b)
		case "/":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			if b == 0 {
				return 0, errors.New("division by zero")
			}
			stack = stack[:len(stack)-2]
			stack = append(stack, a/b)
		default:
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token)
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}