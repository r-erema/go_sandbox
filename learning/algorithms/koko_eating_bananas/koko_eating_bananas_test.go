package kokoeatingbananas_test

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinEatingSpeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		piles []int
		hours int
		want  int
	}{
		{
			name:  "Normal piles",
			piles: []int{30, 11, 23, 4, 20},
			hours: 5,
			want:  30,
		},
		{
			name: "Big piles",
			piles: []int{
				873375536, 395271806, 617254718, 970525912, 634754347, 824202576, 694181619, 20191396, 886462834, 442389139,
				572655464, 438946009, 791566709, 776244944, 694340852, 419438893, 784015530, 588954527, 282060288, 269101141,
				499386849, 846936808, 92389214, 385055341, 56742915, 803341674, 837907634, 728867715, 20958651, 167651719,
				345626668, 701905050, 932332403, 572486583, 603363649, 967330688, 484233747, 859566856, 446838995, 375409782,
				220949961, 72860128, 998899684, 615754807, 383344277, 36322529, 154308670, 335291837, 927055440, 28020467,
				558059248, 999492426, 991026255, 30205761, 884639109, 61689648, 742973721, 395173120, 38459914, 705636911,
				30019578, 968014413, 126489328, 738983100, 793184186, 871576545, 768870427, 955396670, 328003949, 786890382,
				450361695, 994581348, 158169007, 309034664, 388541713, 142633427, 390169457, 161995664, 906356894, 379954831,
				448138536,
			},
			hours: 943223529,
			want:  46,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, minEatingSpeed(tt.piles, tt.hours))
		})
	}
}

// Time O(NlogM), other than the sort invocation, we do a simple linear scan of the list,
// Space O(1), where N is the length of the input array piles and M is the maximum number of bananas in a pile.
func minEatingSpeed(piles []int, hours int) int {
	minSpeed := slices.Max(piles)
	left, right := 1, minSpeed

	eatingTime := func(speed int) int {
		var time int
		for _, pile := range piles {
			time += int(math.Ceil(float64(pile) / float64(speed)))
		}

		return time
	}

	for left <= right {
		potentialMinSpeed := (left + right) / 2

		if eatingTime(potentialMinSpeed) <= hours {
			minSpeed = potentialMinSpeed
			right = potentialMinSpeed - 1
		} else {
			left = potentialMinSpeed + 1
		}
	}

	return minSpeed
}
