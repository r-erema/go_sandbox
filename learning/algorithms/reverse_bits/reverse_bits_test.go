package reversebits_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseBits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number uint32
		want   uint32
	}{
		{
			name:   "bits sequence 1",
			number: 0b00000010100101000001111010011100,
			want:   0b00111001011110000010100101000000,
		},
		{
			name:   "bits sequence 2",
			number: 0b11111111111111111111111111111101,
			want:   0b10111111111111111111111111111111,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, reverseBits(testCase.number))
		})
	}
}

// Time O(N), since we iterate input one time
// Space O(N), since we involve additional allocation to collect reversed bits.
func reverseBits(num uint32) uint32 {
	var res uint32

	for i := 0; i < 32; i++ {
		bit := (num >> i) & 1
		res |= bit << (31 - i)
	}

	return res
}
