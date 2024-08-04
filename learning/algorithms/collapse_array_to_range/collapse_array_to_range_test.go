package collapsearraytorange_test

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollapseArrayToRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  []int
		want string
	}{
		{
			name: "Case 0",
			arr:  []int{1, 4, 5, 2, 3, 9, 8, 11, 0},
			want: "0-5,8-9,11",
		},
		{
			name: "Case 1",
			arr:  []int{1, 3, 2},
			want: "1-3",
		},
		{
			name: "Case 2",
			arr:  []int{1, 4},
			want: "1,4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, collapseArrayToRange(tt.arr))
		})
	}
}

/*
Average, Worst: O(n log n) time | O(n log n) space.
*/
func collapseArrayToRange(arr []int) string {
	sort.Ints(arr)

	ranges := make([]string, 0)
	arrLength := len(arr)
	indexStartRange := 0

	for i := range arrLength {
		if i+1 == arrLength || isGap(float64(arr[i]), float64(arr[i+1])) {
			ranges = append(ranges, buildRange(arr[indexStartRange:i+1]))
			indexStartRange = i + 1
		}
	}

	return strings.Join(ranges, ",")
}

func isGap(number1, number2 float64) bool {
	return math.Abs(number2-number1) > 1
}

func buildRange(arr []int) string {
	length := len(arr)
	if length == 0 {
		return ""
	}

	if length == 1 {
		return strconv.Itoa(arr[0])
	}

	return strconv.Itoa(arr[0]) + "-" + strconv.Itoa(arr[length-1])
}
