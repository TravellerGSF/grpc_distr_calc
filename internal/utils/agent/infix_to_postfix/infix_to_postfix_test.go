package infix_to_postfix

import "testing"

func TestToPostfix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2 + 2", "2 2 +"},
		{"3 + 4 * 2", "3 4 2 * +"},
		{"(1 + 2) * 3", "1 2 + 3 *"},
		{"-5 + 3", "-5 3 +"},
		{"(10 / (5 - 3)) + 1", "10 5 3 - / 1 +"},
	}

	for _, tt := range tests {
		got := ToPostfix(tt.input)
		if got != tt.expected {
			t.Errorf("ToPostfix(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
