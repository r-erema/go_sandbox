package longest_consecutive_sequence

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestLongestConsecutiveSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want int
	}{
		{
			name: "Normal sequence",
			arr:  []int{101, 4, 1, 3, 100, 2},
			want: 4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, longestConsecutive(tt.arr))
		})
	}
}

// Time: O(n), since we iterate only through the chain of the sequential numbers thanks to the map
// TimeL O(n), since we involve only a map which equals an input
func longestConsecutive(input []int) int {
	numsMap := make(map[int]struct{}, len(input))
	for _, num := range input {
		numsMap[num] = struct{}{}
	}

	longest := 0
	for n := range numsMap {
		if _, ok := numsMap[n-1]; !ok {
			length := 0
			for {
				_, nextExists := numsMap[n+length]
				if !nextExists {
					break
				}

				length++
				longest = max(longest, length)
			}
		}
	}

	return longest
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}

	return n2
}
