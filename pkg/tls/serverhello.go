package tls

import (
	"encoding/binary"
	"fmt"
)

func encodeServerHello(key publicKey, secret [32]byte, clientCipher cipherSuite) ([]byte, error) {
	helloMessage := []byte{
		// Record header
		0x16,       // type is 0x16 (handshake record)
		0x03, 0x04, // protocol version is "3.4" (also known as TLS 1.3)
	}

	var handshakeHeader []byte

	handshakeHeader = append(handshakeHeader, serverType[:]...)

	handshakeData := []byte{
		// Server Version
		0x03, 0x04, // protocol version is "3.4" (also known as TLS 1.3)
	}

	// Server random
	handshakeData = append(handshakeData, secret[:]...)

	// Session ID
	handshakeData = append(handshakeData, []byte{
		// No Session ID
		0x00,
	}...)

	// Cipher Suite
	handshakeData = append(handshakeData, clientCipher[:]...)

	// Compression Method
	handshakeData = append(handshakeData, []byte{0x00}...)

	encodedExtensions := serverExtensions(key)
	handshakeData = binary.BigEndian.AppendUint16(handshakeData, uint16(len(encodedExtensions)))
	handshakeData = append(handshakeData, encodedExtensions...)

	handshakeHeader, err := BigEndianAppend24(handshakeHeader, uint32(len(handshakeData)))
	if err != nil {
		return nil, fmt.Errorf("failed to append handshake header length: %w", err)
	}

	helloMessage = binary.BigEndian.AppendUint16(helloMessage, uint16(len(handshakeHeader)+len(handshakeData)))
	helloMessage = append(helloMessage, handshakeHeader...)
	helloMessage = append(helloMessage, handshakeData...)

	return helloMessage, nil
}
