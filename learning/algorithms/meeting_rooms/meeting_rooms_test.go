package meetingrooms_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeetingRooms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		meetings [][2]int
		want     bool
	}{
		{
			name:     "can not attend all meetings",
			meetings: [][2]int{{0, 30}, {5, 10}, {15, 20}},
			want:     false,
		},
		{
			name:     "can attend all meetings",
			meetings: [][2]int{{5, 8}, {9, 15}},
			want:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, canAttendMeetings(tt.meetings))
		})
	}
}

// Time O(NlogN), other than the sort invocation, we do a simple linear scan of the list,
// Space O(logN), due to the sorting.
func canAttendMeetings(meetings [][2]int) bool {
	sort.Slice(meetings, func(i, j int) bool {
		return meetings[i][0] < meetings[j][0]
	})

	for i := 1; i < len(meetings); i++ {
		if meetings[i][0] <= meetings[i-1][1] {
			return false
		}
	}

	return true
}
