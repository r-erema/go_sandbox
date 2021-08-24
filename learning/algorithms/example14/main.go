package example14

const (
	alphabetLength         = 26
	alphabetStartCharPoint = 96
	lastAlphabetChar       = 122
)

/*
	Average, Worst: O(n) time | O(1) space
*/
func CaesarCipherEncrypt(str string, shift int32) string {
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
