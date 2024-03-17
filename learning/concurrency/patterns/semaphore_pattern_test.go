package patterns_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderProcessor struct {
	semaphore chan struct{}
	wg        *sync.WaitGroup
}

func (p *orderProcessor) process(_ int) {
	p.semaphore <- struct{}{}
	defer func() {
		<-p.semaphore
		p.wg.Done()
	}()

	time.Sleep(time.Millisecond * 100)
}

func TestProcess(t *testing.T) {
	t.Parallel()

	const semaphoreSize = 2

	tests := []struct {
		name                              string
		orderIDs                          []int
		expectedConcurrentProcessingCount int
	}{
		{
			name:                              "1 order",
			orderIDs:                          []int{1},
			expectedConcurrentProcessingCount: 1,
		},
		{
			name:                              "2 orders",
			orderIDs:                          []int{1, 2},
			expectedConcurrentProcessingCount: 2,
		},
		{
			name:                              "3 order2",
			orderIDs:                          []int{1, 2, 3},
			expectedConcurrentProcessingCount: 2,
		},
		{
			name:                              "10 orders",
			orderIDs:                          []int{1, 2, 3, 11, 22, 33, 4444, 5, 9, 10},
			expectedConcurrentProcessingCount: 2,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			processor := orderProcessor{
				semaphore: make(chan struct{}, semaphoreSize),
				wg:        new(sync.WaitGroup),
			}

			for _, id := range testCase.orderIDs {
				processor.wg.Add(1)
				go processor.process(id)
			}

			for len(processor.semaphore) < testCase.expectedConcurrentProcessingCount {
				time.Sleep(time.Millisecond)
			}

			assert.Equal(t, testCase.expectedConcurrentProcessingCount, len(processor.semaphore))

			processor.wg.Wait()
			assert.Equal(t, 0, len(processor.semaphore))
		})
	}
}

var (
	errAllLicensesInUse    = errors.New("all licenses are in use")
	errNoLicensesToRelease = errors.New("no license to release")
)

type licenseManager struct {
	semaphore                 chan struct{}
	mu                        sync.Mutex
	maxLicenses, usedLicenses int
}

func buildLicenseManager(maxLicenses int) *licenseManager {
	return &licenseManager{
		semaphore:   make(chan struct{}, maxLicenses), // Define the semaphore max number
		maxLicenses: maxLicenses,
	}
}

func (lm *licenseManager) Acquire() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lm.usedLicenses >= lm.maxLicenses {
		return errAllLicensesInUse
	}

	lm.semaphore <- struct{}{}
	lm.usedLicenses++

	return nil
}

func (lm *licenseManager) Release() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	select {
	case <-lm.semaphore:
		lm.usedLicenses--

		return nil
	default:
		return errNoLicensesToRelease
	}
}

func (lm *licenseManager) UsedLicenses() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	return lm.usedLicenses
}

func TestLicenseManager(t *testing.T) {
	t.Parallel()

	t.Run("Acquire", func(t *testing.T) {
		t.Parallel()

		licMgr := buildLicenseManager(2)
		err := licMgr.Acquire()
		assert.NoError(t, err)
		assert.Equal(t, 1, licMgr.UsedLicenses())
	})

	t.Run("Release", func(t *testing.T) {
		t.Parallel()

		licMgr := buildLicenseManager(2)
		err := licMgr.Acquire()
		require.NoError(t, err)

		err = licMgr.Release()
		assert.NoError(t, err)
		assert.Equal(t, 0, licMgr.UsedLicenses())
	})

	t.Run("Exceed", func(t *testing.T) {
		t.Parallel()

		licMgr := buildLicenseManager(1)
		err := licMgr.Acquire()
		assert.NoError(t, err)
		err = licMgr.Acquire()
		assert.Error(t, err)
	})

	t.Run("Release Without Acquire", func(t *testing.T) {
		t.Parallel()

		licMgr := buildLicenseManager(1)
		err := licMgr.Release()
		assert.Error(t, err)
	})
}
