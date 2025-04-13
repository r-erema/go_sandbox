package carfleet_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCarFleet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		target   int
		position []int
		speed    []int
		want     int
	}{
		{
			name:     "3 fleets of 5 cars",
			target:   12,
			position: []int{10, 8, 0, 5, 3},
			speed:    []int{2, 4, 1, 1, 3},
			want:     3,
		},
		{
			name:     "1 fleets of 1 car",
			target:   10,
			position: []int{3},
			speed:    []int{3},
			want:     1,
		},
		{
			name:     "1 fleets of 3 cars",
			target:   10,
			position: []int{0, 4, 2},
			speed:    []int{2, 1, 3},
			want:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, carFleet(tt.target, tt.position, tt.speed))
		})
	}
}

func carFleet(target int, position, speed []int) int {
	cars := make([][2]int, len(position))
	for i := range position {
		cars[i] = [2]int{position[i], speed[i]}
	}

	sort.Slice(cars, func(i, j int) bool {
		return cars[i][0] < cars[j][0]
	})

	fleets := len(cars)

	currentCarTimeToDestination := float32(
		target-cars[len(cars)-1][0],
	) / float32(
		cars[len(cars)-1][1],
	)

	for i := len(cars) - 1; i > 0; i-- {
		previousCarTimeToDestination := float32(target-cars[i-1][0]) / float32(cars[i-1][1])
		if previousCarTimeToDestination <= currentCarTimeToDestination {
			fleets--
		} else {
			currentCarTimeToDestination = previousCarTimeToDestination
		}
	}

	return fleets
}
