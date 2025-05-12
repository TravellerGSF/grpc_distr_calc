package validator

import "testing"

func TestIsValidExpression(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"2 + 2", true},
		{"(3 + 4) * 5", true},
		{"", false},
		{"2 + (3 *", false},
		{"2 + abc", false},
		{"3 + 4)", false},
		{"(3 + 4", false},
	}

	for _, tt := range tests {
		if got := IsValidExpression(tt.input); got != tt.valid {
			t.Errorf("IsValidExpression(%q) = %v; want %v", tt.input, got, tt.valid)
		}
	}
}
