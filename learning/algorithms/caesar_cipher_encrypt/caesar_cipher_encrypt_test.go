package caesarcipherencrypt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaesarCipherEncrypt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, str, want string
		shift           int32
	}{
		{
			name:  "Case 0",
			str:   "xyz",
			shift: 2,
			want:  "zab",
		},
		{
			name:  "Case 1",
			str:   "lorem ipsum dolor sit amet, consectetur adipiscing elit",
			shift: 732,
			want:  "psviq mtwyq hspsv wmx eqix, gsrwigxixyv ehmtmwgmrk ipmx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, caesarCipherEncrypt(tt.str, tt.shift))
		})
	}
}

const (
	alphabetLength         = 26
	alphabetStartCharPoint = 96
	lastAlphabetChar       = 122
)

/*
Average, Worst: O(n) time | O(1) space.
*/
func caesarCipherEncrypt(str string, shift int32) string {
	var result, encryptedSymbol string

	isLetter := func(char int32) bool {
		return char > alphabetStartCharPoint && char <= lastAlphabetChar
	}

	shift %= alphabetLength

	for _, char := range str {
		if !isLetter(char) {
			encryptedSymbol = string(char)
		} else {
			nextChar := char + shift
			if nextChar > lastAlphabetChar {
				encryptedSymbol = string(alphabetStartCharPoint + nextChar - lastAlphabetChar)
			} else {
				encryptedSymbol = string(nextChar)
			}
		}

		result += encryptedSymbol
	}

	return result
}
