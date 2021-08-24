package example15

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

/*
	Average, Worst: O(n log n) time | O(n log n) space
*/
func CollapseArrayToRange(arr []int) string {
	sort.Ints(arr)

	ranges := make([]string, 0)
	arrLength := len(arr)
	indexStartRange := 0

	for i := 0; i < arrLength; i++ {
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
