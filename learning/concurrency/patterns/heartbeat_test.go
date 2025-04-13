package patterns_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeartbeat(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})

	const timeout = time.Millisecond * 100
	heartbeat, results := doWork(done, timeout/2)

	time.AfterFunc(time.Millisecond*500, func() {
		close(done)
	})

	var heartbeatsCount, resultsCount int

	func() {
		for {
			select {
			case _, ok := <-heartbeat:
				if !ok {
					return
				}

				heartbeatsCount++
			case res, ok := <-results:
				if !ok {
					return
				}

				require.Contains(t, res, "some work is done")

				resultsCount++
			}
		}
	}()

	assert.Positive(t, heartbeatsCount)
	assert.Positive(t, resultsCount)
}

//nolint:nonamedreturns
func doWork(done <-chan struct{}, pulseInterval time.Duration) (
	heartbeat chan struct{}, results chan string,
) {
	heartbeat = make(chan struct{})
	results = make(chan string)
	pulse := time.Tick(pulseInterval)
	workGen := time.Tick(2 * pulseInterval)
	sendResult := func(resultTime time.Time) {
		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse(heartbeat)
			case results <- fmt.Sprintf("some work is done: %s", resultTime):
				return
			}
		}
	}

	go func() {
		defer close(heartbeat)
		defer close(results)

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse(heartbeat)
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

func sendPulse(heartbeat chan<- struct{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}
