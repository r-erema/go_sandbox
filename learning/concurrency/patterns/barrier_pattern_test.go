package patterns_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"
)

type barrier struct {
	totalGoroutines,
	goroutinesReachedBarrier int
	mu   sync.Mutex
	cond *sync.Cond
}

func newBarrier(goroutinesCountGoingToBarrier int) *barrier {
	barr := &barrier{
		totalGoroutines:          goroutinesCountGoingToBarrier,
		goroutinesReachedBarrier: 0,
	}
	barr.cond = sync.NewCond(&barr.mu)

	return barr
}

func (b *barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.goroutinesReachedBarrier++
	if b.goroutinesReachedBarrier >= b.totalGoroutines {
		b.goroutinesReachedBarrier = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
}

func TestBarrierFirstInstance(t *testing.T) {
	t.Parallel()

	var (
		wg     sync.WaitGroup
		mutex  sync.Mutex
		sum    int
		amount = 5
	)

	barr := newBarrier(5)

	rand.Seed(time.Now().UnixNano())

	for range amount {
		wg.Add(1)

		go func() {
			defer wg.Done()

			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			barr.Wait()

			mutex.Lock()
			defer mutex.Unlock()

			sum++
		}()
	}

	wg.Wait()

	assert.Equal(t, amount, sum)
}

func makeRequest(url string, barrier *barrier, respBodyCh chan<- []byte, errCh chan<- error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		errCh <- fmt.Errorf("error creating request to %s: %w", url, err)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("error making request to %s: %w", url, err)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- fmt.Errorf("error reading response body: %w", err)

		return
	}

	err = resp.Body.Close()
	if err != nil {
		errCh <- fmt.Errorf("error body close for %s: %w", url, err)
	}

	close(errCh)

	barrier.Wait()

	respBodyCh <- body
}

func TestAllHTTPRequestsShouldBeOKOtherwiseReturnError(t *testing.T) {
	t.Parallel()

	rand.Seed(time.Now().UnixNano())

	microserviceDepositRate := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		_, err := w.Write([]byte("0.02"))
		assert.NoError(t, err)
	}))

	microserviceUserAccountBalance := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		_, err := w.Write([]byte("500"))
		w.WriteHeader(http.StatusInternalServerError)
		assert.NoError(t, err)
	}))

	tests := []struct {
		name     string
		assertFn func(
			barrier *barrier,
			respBodyDepositRateCh,
			respBodyUserBalanceRateCh chan []byte,
			respErrDepositRateCh,
			respErrUserBalanceRateCh chan error)
	}{
		{
			name: "Both requests OK",
			assertFn: func(barrier *barrier, bodyDepositCh, bodyBalanceCh chan []byte, errDepositCh, errBalanceCh chan error) {
				go makeRequest(microserviceDepositRate.URL, barrier, bodyDepositCh, errDepositCh)
				go makeRequest(microserviceUserAccountBalance.URL, barrier, bodyBalanceCh, errBalanceCh)

				require.NoError(t, <-errDepositCh)
				require.NoError(t, <-errBalanceCh)

				rate, err := strconv.ParseFloat(string(<-bodyDepositCh), 64)
				require.NoError(t, err)
				balance, err := strconv.ParseFloat(string(<-bodyBalanceCh), 64)
				require.NoError(t, err)

				assert.InEpsilon(t, float64(10), rate*balance, 0)
			},
		},
		{
			name: "One request failed",
			assertFn: func(barrier *barrier, bodyDepositCh, bodyBalanceCh chan []byte, errDepositCh, errBalanceCh chan error) {
				go makeRequest("microserviceDepositRate.bad.url", barrier, bodyDepositCh, errDepositCh)
				go makeRequest(microserviceUserAccountBalance.URL, barrier, bodyBalanceCh, errBalanceCh)

				require.Error(t, <-errDepositCh)
				require.NoError(t, <-errBalanceCh)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.assertFn(newBarrier(2), make(chan []byte), make(chan []byte), make(chan error), make(chan error))
		})
	}
}
