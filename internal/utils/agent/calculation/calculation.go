package calculation

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	operatorPlusDelay     time.Duration
	operatorMinusDelay    time.Duration
	operatorMultiplyDelay time.Duration
	operatorDivideDelay   time.Duration
)

func init() {
	operatorPlusDelay = readDurationEnv("OPERATOR_PLUS_DELAY", 1000)
	operatorMinusDelay = readDurationEnv("OPERATOR_MINUS_DELAY", 1000)
	operatorMultiplyDelay = readDurationEnv("OPERATOR_MULTIPLY_DELAY", 7000)
	operatorDivideDelay = readDurationEnv("OPERATOR_DIVIDE_DELAY", 12000)
}

func readDurationEnv(key string, defaultValue int) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return time.Duration(defaultValue) * time.Millisecond
	}
	duration, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("Invalid value for %s: %s, using default %d ms\n", key, val, defaultValue)
		return time.Duration(defaultValue) * time.Millisecond
	}
	return time.Duration(duration) * time.Millisecond
}

func Evaluate(expr string) (float64, error) {
	var stack Stack
	tokens := strings.Split(expr, " ")
	for _, token := range tokens {
		if token == "+" || token == "-" || token == "*" || token == "/" {
			op1 := stack.Pop()
			op2 := stack.Pop()
			ans, err := Calculate(op1, op2, token)
			if err != nil {
				return 0, err
			}
			stack.Push(ans)
		} else {
			op, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("неверный формат числа: %s", token)
			}
			stack.Push(op)
		}
	}
	return stack.Pop(), nil
}

func Calculate(op1, op2 float64, operator string) (float64, error) {
	switch operator {
	case "+":
		time.Sleep(operatorPlusDelay)
		return op2 + op1, nil
	case "-":
		time.Sleep(operatorMinusDelay)
		return op2 - op1, nil
	case "*":
		time.Sleep(operatorMultiplyDelay)
		return op2 * op1, nil
	case "/":
		time.Sleep(operatorDivideDelay)
		if op1 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return op2 / op1, nil
	default:
		return 0, fmt.Errorf("неизвестная операция")
	}
}
