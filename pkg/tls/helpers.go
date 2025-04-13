package tls

import (
	"crypto/rand"
	"fmt"
)

func Rand32Bytes() ([secretLength]byte, error) {
	buf := make([]byte, secretLength)
	if _, err := rand.Read(buf); err != nil {
		return [32]byte{}, fmt.Errorf("rand Read error: %w", err)
	}

	return [32]byte(buf), nil
}

func BigEndianAppend24(buf []byte, val uint32) ([]byte, error) {
	var maxBits uint32 = 0xffffff
	if val > maxBits {
		return nil, NewError(fmt.Sprintf("value out of range of 24 bits: %d", val))
	}

	offset16, offset8 := 16, 8

	buf = append(buf, byte(val>>offset16), byte(val>>offset8), byte(val))

	return buf, nil
}
