package permutationinstring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInclusion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input1, input2 string
		want           bool
	}{
		{
			name:   "permutation exists",
			input1: "a",
			input2: "ab",
			want:   true,
		},
		{
			name:   "permutation does not exist",
			input1: "ab",
			input2: "eidboaoo",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, checkInclusion(tt.input1, tt.input2))
		})
	}
}

func checkInclusion(str1, str2 string) bool {
	if len(str1) > len(str2) {
		return false
	}

	var frequencyStr1, frequencyStr2 [26]rune

	for i := range str1 {
		frequencyStr1[str1[i]-'a']++
		frequencyStr2[str2[i]-'a']++
	}

	refreshFrequencyS2 := func(left, right int) {
		if left > 0 {
			frequencyStr2[str2[left-1]-'a']--
		}

		if right < len(str2) {
			frequencyStr2[str2[right]-'a']++
		}
	}

	left, right := 1, len(str1)
	for right < len(str2) {
		if frequencyStr1 == frequencyStr2 {
			return true
		}

		refreshFrequencyS2(left, right)

		left++
		right++
	}

	return frequencyStr1 == frequencyStr2
}
