package infix_to_postfix

import (
	"strings"
)

func ToPostfix(s string) string {
	var stack Stack
	postfix := ""
	length := len(s)

	s = strings.ReplaceAll(s, ",", ".")

	for i := 0; i < length; i++ {
		char := string(s[i])
		if char == " " {
			continue
		}

		if char == "(" {
			stack.Push(char)
		} else if char == ")" {
			for !stack.Empty() {
				str, _ := stack.TopFunc().(string)
				if str == "(" {
					break
				}
				postfix += " " + str
				stack.Pop()
			}
			stack.Pop()
		} else if char == "-" && (i == 0 || s[i-1] == '(') {
			number := "-"
			j := i + 1
			for ; j < length && IsOperand(s[j]); j++ {
				number += string(s[j])
			}
			postfix += " " + number
			i = j - 1
		} else if !IsOperator(s[i]) {
			j := i
			number := ""
			for ; j < length && (IsOperand(s[j]) || s[j] == '.'); j++ {
				number += string(s[j])
			}
			postfix += " " + number
			i = j - 1
		} else {
			for !stack.Empty() {
				top, _ := stack.TopFunc().(string)
				if top == "(" || !HasHigherPrecedence(top, char) {
					break
				}
				postfix += " " + top
				stack.Pop()
			}
			stack.Push(char)
		}
	}
	for !stack.Empty() {
		str, _ := stack.Pop().(string)
		postfix += " " + str
	}
	return strings.TrimSpace(postfix)
}

func IsOperator(c uint8) bool {
	return strings.ContainsAny(string(c), "+-*/")
}

func IsOperand(c uint8) bool {
	return c >= '0' && c <= '9' || c == '.'
}

func Precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return -1
}

func HasHigherPrecedence(op1, op2 string) bool {
	op1Prec := Precedence(op1)
	op2Prec := Precedence(op2)
	return op1Prec >= op2Prec
}
