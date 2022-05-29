package utils

import (
	"bytes"
	"fmt"
	"sync"
)

type ThreadSafeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (tsb *ThreadSafeBuffer) Write(p []byte) (int, error) {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()

	bytesCount, err := tsb.buffer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("internal buffer writing error %w", err)
	}

	return bytesCount, nil
}

func (tsb *ThreadSafeBuffer) String() string {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()

	return tsb.buffer.String()
}
