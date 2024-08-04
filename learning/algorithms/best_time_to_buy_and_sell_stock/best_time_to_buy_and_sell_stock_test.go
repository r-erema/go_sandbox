package besttimetobuyandsellstock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestTimeToBuyAndSellStock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{
			name:   "Simple set",
			prices: []int{7, 1, 5, 3, 6, 4},
			want:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, maxProfit(tt.prices))
		})
	}
}

// Time O(n), since we should iterate all the input
// Space O(1), sine we don't allocate any additional memory.
func maxProfit(prices []int) int {
	buyIndex, sellIndex := 0, 0
	profit := 0

	for sellIndex < len(prices) {
		if prices[buyIndex] > prices[sellIndex] {
			buyIndex = sellIndex
		}

		profit = max(profit, prices[sellIndex]-prices[buyIndex])

		sellIndex++
	}

	return profit
}
