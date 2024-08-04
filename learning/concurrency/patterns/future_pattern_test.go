package patterns_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errEmptyData = errors.New("data slice is empty")

type future struct {
	result chan float64
	err    chan error
}

func (f *future) Get() (float64, error) {
	select {
	case res := <-f.result:
		return res, nil
	case err := <-f.err:
		return 0, err
	}
}

func calculateAverage(data []int) *future {
	promise := &future{
		result: make(chan float64),
		err:    make(chan error),
	}

	go func() {
		if len(data) == 0 {
			promise.err <- errEmptyData

			return
		}

		var sum int

		for _, datum := range data {
			time.Sleep(time.Millisecond * 100)

			sum += datum
		}

		promise.result <- float64(sum) / float64(len(data))
	}()

	return promise
}

func TestCalculateAverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []int
		expected float64
		hasError bool
	}{
		{
			name:     "Average of non-empty slice",
			data:     []int{10, 20, 30, 40, 50},
			expected: 30,
			hasError: false,
		},
		{
			name:     "Average of empty slice",
			data:     []int{},
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := calculateAverage(tt.data)
			avg, err := f.Get()

			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.InEpsilon(t, tt.expected, avg, 0)
		})
	}
}

type result struct {
	body []byte
	err  error
}

type futureFetch struct {
	result chan result
}

func (f *futureFetch) Get() result {
	return <-f.result
}

func newFutureFetch(f func() (result, error)) *futureFetch {
	res := make(chan result, 1)

	go func() {
		r, err := f()
		res <- result{body: r.body, err: err}
		close(res)
	}()

	return &futureFetch{result: res}
}

func fetch(url string) (result, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return result{}, fmt.Errorf("creation request error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result{}, fmt.Errorf("GET request error: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		log.Print(err)
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result{}, fmt.Errorf("read body error: %w", err)
	}

	return result{body: body}, nil
}

func TestFetch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Millisecond * 100)

		_, err := writer.Write([]byte("fetched data"))
		assert.NoError(t, err)
		writer.WriteHeader(http.StatusOK)
	}))

	f := newFutureFetch(func() (result, error) {
		return fetch(server.URL)
	})

	res := f.Get()
	require.NoError(t, res.err)
	assert.Equal(t, []byte("fetched data"), res.body)
}
