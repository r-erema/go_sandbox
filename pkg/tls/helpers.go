package tls

import (
	"crypto/rand"
	"fmt"
)

func Rand32Bytes() ([32]byte, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return [32]byte{}, fmt.Errorf("rand Read error: %w", err)
	}
	return [32]byte(buf), nil
}

func BigEndianAppend24(b []byte, v uint32) ([]byte, error) {
	if v > 0xffffff {
		return nil, fmt.Errorf("value out of range of 24 bits: %d", v)
	}

	b = append(b, byte(v>>16), byte(v>>8), byte(v))

	return b, nil
}
