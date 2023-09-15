package validanagram_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidAnagram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, word1, word2 string
		want               bool
	}{
		{
			name:  "anagram",
			word1: "anagram",
			word2: "nagaram",
			want:  true,
		},
		{
			name:  "not anagram",
			word1: "rat",
			word2: "car",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isAnagram(tt.word1, tt.word2))
		})
	}
}

// Time O(n), since we need to iterate each element in a word
// Space O(1) - we use an array with length 26 to be able to store all 26 ASCII characters,
//
//	where indexes are difference of ASCII code of lower case letter and minimum possible code of lower case letter,
//	i.e. 'a' == 97, e.g. if we have letter 'z', it's code 122, but we can't add it into 26 length array,
//	so we subtract 122-97=25 and can put it into array, i.e. array[25]++
func isAnagram(word1, word2 string) bool {
	if len(word1) != len(word2) {
		return false
	}

	arr := [26]int{}

	for i := 0; i < len(word1); i++ {
		arr[word1[i]-'a']++
		arr[word2[i]-'a']--
	}

	return arr == [26]int{}
}
