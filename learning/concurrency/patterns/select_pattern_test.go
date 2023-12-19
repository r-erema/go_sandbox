package patterns_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelectPatternNonBlockingCases(t *testing.T) {
	t.Parallel()

	messages := make(chan string)

	tests := []struct {
		name       string
		assertFunc func(t *testing.T)
	}{
		{
			name: "Non-blocking send",
			assertFunc: func(t *testing.T) {
				t.Helper()
				select {
				case messages <- "test":
					t.Fail()
				default:
					assert.True(t, true)
				}
			},
		},
		{
			name: "Non-blocking receive",
			assertFunc: func(t *testing.T) {
				t.Helper()
				select {
				case <-messages:
					t.Fail()
				default:
					assert.True(t, true)
				}
			},
		},
		{
			name: "Non-blocking multiway select",
			assertFunc: func(t *testing.T) {
				t.Helper()
				select {
				case <-messages:
					t.Fail()
				case messages <- "test":
					t.Fail()
				default:
					assert.True(t, true)
				}
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.assertFunc(t)
		})
	}
}

func TestSelectPatternTimeoutCases(t *testing.T) {
	t.Parallel()

	const mockFetchedDataValue = "data fetched"

	mockFetchData := func(response chan<- string, delay time.Duration) {
		time.Sleep(delay)
		response <- mockFetchedDataValue
	}

	tests := []struct {
		name       string
		assertFunc func(t *testing.T)
	}{
		{
			name: "Timeout",
			assertFunc: func(t *testing.T) {
				t.Helper()
				dataCh := make(chan string)
				delay := time.Millisecond * 100
				go mockFetchData(dataCh, delay)
				select {
				case <-dataCh:
					t.Fail()
				case <-time.After(delay - time.Millisecond):
					assert.True(t, true)
				}
			},
		},
		{
			name: "Fetch data before timeout",
			assertFunc: func(t *testing.T) {
				t.Helper()
				dataCh := make(chan string)
				delay := time.Millisecond * 100
				go mockFetchData(dataCh, delay)
				select {
				case data := <-dataCh:
					assert.Equal(t, mockFetchedDataValue, data)
				case <-time.After(delay + time.Millisecond):
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.assertFunc(t)
		})
	}
}
