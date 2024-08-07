package type_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringIteration(t *testing.T) {
	t.Parallel()

	str := "abcАБВ🙂"

	assert.Len(t, str, 13)

	iterationsCount := 0

	for i := range len(str) {
		assert.IsType(t, byte(0), str[i])

		iterationsCount++
	}

	assert.Equal(t, 13, iterationsCount)

	iterationsCount = 0

	for i, r := range str {
		assert.Contains(t, [7]int{0, 1, 2, 3, 5, 7, 9}, i)
		assert.IsType(t, '1', r)

		iterationsCount++
	}

	assert.Equal(t, 7, iterationsCount)
}
