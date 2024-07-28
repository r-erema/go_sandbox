package timebasedkeyvaluestore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeBasedKeyValueStore(t *testing.T) {
	t.Parallel()

	timeMap := Constructor()
	timeMap.Set("foo", "bar", 1)
	assert.Equal(t, "bar", timeMap.Get("foo", 1))
	assert.Equal(t, "bar", timeMap.Get("foo", 3))
	timeMap.Set("foo", "bar2", 4)
	assert.Equal(t, "bar2", timeMap.Get("foo", 4))
	assert.Equal(t, "bar2", timeMap.Get("foo", 5))
}

type timestampValue struct {
	timestamp int
	value     string
}

type TimeMap struct {
	timeSeriesMap map[string][]timestampValue
}

func Constructor() TimeMap {
	return TimeMap{
		timeSeriesMap: make(map[string][]timestampValue),
	}
}

func (tm *TimeMap) Set(key, value string, timestamp int) {
	if tm.timeSeriesMap[key] == nil {
		tm.timeSeriesMap[key] = make([]timestampValue, 0)
	}

	tm.timeSeriesMap[key] = append(tm.timeSeriesMap[key], timestampValue{value: value, timestamp: timestamp})
}

func (tm *TimeMap) Get(key string, timestamp int) string {
	var res string

	if timestampValues, ok := tm.timeSeriesMap[key]; ok {
		left, right := 0, len(timestampValues)-1
		for left <= right {
			mid := (left + right) / 2

			if timestamp >= timestampValues[mid].timestamp {
				res = timestampValues[mid].value
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return res
}
