package validator

var validChars = map[rune]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	'+': true, '-': true, '*': true, '/': true,
	'(': true, ')': true, ' ': true, ',': true,
	'.': true,
}

func IsValidExpression(expression string) bool {
	if !isValidCharacters(expression) || !isBalanced(expression) {
		return false
	}
	return true
}

func isValidCharacters(expression string) bool {
	if expression == "" {
		return false
	}
	for _, char := range expression {
		if !validChars[char] {
			return false
		}
	}
	return true
}

func isBalanced(expression string) bool {
	count := 0
	for _, char := range expression {
		if char == '(' {
			count++
		} else if char == ')' {
			count--
			if count < 0 {
				return false
			}
		}
	}
	return count == 0
}
