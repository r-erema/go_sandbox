package encodeanddecodestrings_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecodeStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		source     []string
		encodedStr string
	}{
		{
			name:       "Simple words",
			source:     []string{"lint", "code", "love", "you"},
			encodedStr: "4#lint4#code4#love3#you",
		},
		{
			name:       "Long words",
			source:     []string{"estimation#", "#highlighted", "br#ing_it_to_us"},
			encodedStr: "11#estimation#12##highlighted15#br#ing_it_to_us",
		},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			encoded := encode(testCase.source)
			assert.Equal(t, testCase.encodedStr, encoded)
			decoded := decode(encoded)
			assert.Equal(t, testCase.source, decoded)
		})
	}
}

// Time O(N) where N is each string
// Time (1), since we don't use any additional space.
func encode(input []string) string {
	var encoded string
	for _, s := range input {
		encoded += fmt.Sprintf("%d#%s", len(s), s)
	}

	return encoded
}

// Time O(N) where N is each string
// Time (N), since we use array to collect sliced words.
func decode(input string) []string {
	var (
		lengthStr, word string
		res             []string
	)

	for len(input) > 0 {
		i := 0
		for input[i] != '#' {
			i++
		}

		lengthStr, input = input[:i], input[i:]

		length, _ := strconv.Atoi(lengthStr)

		word, input = input[1:length+1], input[length+1:]

		res = append(res, word)
	}

	return res
}
