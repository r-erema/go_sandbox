package evaluatereversepolishnotation_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateReversePolishNotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		tokens []string
		want   int
	}{
		{
			name:   "Short tokens sequence",
			tokens: []string{"2", "1", "+", "3", "*"},
			want:   9,
		},
		{
			name:   "Long tokens sequence",
			tokens: []string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"},
			want:   22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, evaluateReversePolishNotation(tt.tokens))
		})
	}
}

// Time O(n), since we iterate each element of input
// Time O(n), since we involve stack equals input.
func evaluateReversePolishNotation(tokens []string) int {
	var (
		stack              []int
		operand1, operand2 int
	)

	for _, token := range tokens {
		switch token {
		case "+":
			operand1, operand2, stack = getAndPopLastOperand(stack)
			stack = append(stack, operand1+operand2)
		case "-":
			operand1, operand2, stack = getAndPopLastOperand(stack)
			stack = append(stack, operand1-operand2)
		case "*":
			operand1, operand2, stack = getAndPopLastOperand(stack)
			stack = append(stack, operand1*operand2)
		case "/":
			operand1, operand2, stack = getAndPopLastOperand(stack)
			stack = append(stack, operand1/operand2)
		default:
			s, _ := strconv.Atoi(token)
			stack = append(stack, s)
		}
	}

	return stack[0]
}

func getAndPopLastOperand(stack []int) (preLast, last int, newStack []int) { //nolint:nonamedreturns
	return stack[len(stack)-2], stack[len(stack)-1], stack[:len(stack)-2]
}
