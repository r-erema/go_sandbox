package isvalidsubsequence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateSubsequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		sequence    []int
		subSequence []int
		isValid     bool
	}{
		{
			name:        "Test case 0",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{1, 6, -1, 10},
			isValid:     true,
		},
		{
			name:        "Test case 1",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{5, 6, -1, 10},
			isValid:     true,
		},
		{
			name:        "Test case 2",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{5, 6, -1, 22},
			isValid:     false,
		},
		{
			name:        "Test case 3",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{4},
			isValid:     false,
		},
		{
			name:        "Test case 4",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{},
			isValid:     true,
		},
		{
			name:        "Test case 5",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{5, 1, 22, 25, 6, -1, 8, 10, 9},
			isValid:     false,
		},
		{
			name:        "Test case 6",
			sequence:    []int{5, 1, 22, 25, 6, -1, 8, 10},
			subSequence: []int{3, 5, 1, 22, 25, 6, -1, 8, 10},
			isValid:     false,
		},
		{
			name:        "Test case 7",
			sequence:    []int{1, 4, 8, 8, 9, 9, 6, 3, 6, 0, 0, 6, 15, 7, 15},
			subSequence: []int{4, 8, 6, 15},
			isValid:     true,
		},
		{
			name:        "Test case 8",
			sequence:    []int{1, 4, 8, 8, 9, 9, 6, 3, 6, 0, 0, 15, 7, 15},
			subSequence: []int{4, 8, 0, 6, 15},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := isValidSubsequence(testCase.subSequence, testCase.sequence)
			assert.Equal(t, testCase.isValid, result)
			result = isValidSubsequence2(testCase.subSequence, testCase.sequence)
			assert.Equal(t, testCase.isValid, result)
		})
	}
}

func isValidSubsequence(subSequence, sequence []int) bool {
	secIdx, subSeqIdx := 0, 0
	sequenceLength, subSequenceLength := len(sequence), len(subSequence)

	for secIdx < sequenceLength && subSeqIdx < subSequenceLength {
		if sequence[secIdx] == subSequence[subSeqIdx] {
			subSeqIdx++
		}
		secIdx++
	}

	return subSeqIdx == len(subSequence)
}

func isValidSubsequence2(subSequence, sequence []int) bool {
	subSeqIdx := 0
	for _, value := range sequence {
		if subSeqIdx == len(subSequence) {
			break
		}

		if value == subSequence[subSeqIdx] {
			subSeqIdx++
		}
	}

	return subSeqIdx == len(subSequence)
}
