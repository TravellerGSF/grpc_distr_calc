package calculation

import (
	"testing"

	"github.com/TravellerGSF/grpc_distr_calc/internal/utils/agent/infix_to_postfix"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"2 + 2", 4},
		{"5 - 3", 2},
		{"4 * 2", 8},
		{"10 / 2", 5},
		{"(1 + 2) * 3", 9},
		{"(10 / 2) + 3", 8},
		{"2 * (3 + 4)", 14},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * (4 + 5)", 45},
		{"2 + 3 * (4 + 5)", 23},
		{"2 * (3 + 4) / 5", 2.8},
	}

	for _, test := range tests {
		postfix := infix_to_postfix.ToPostfix(test.expression)
		result, err := Evaluate(postfix)
		if err != nil {
			t.Errorf("Evaluate(%q) returned an error: %v", test.expression, err)
		}
		if result != test.expected {
			t.Errorf("Evaluate(%q) = %v; want %v", test.expression, result, test.expected)
		}
	}
}
