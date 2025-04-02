package tls

import (
	"encoding/binary"
	"fmt"
)

func encodeClientHello(hostNames []string, keys []publicKey, secret [32]byte) ([]byte, error) {
	helloMessage := []byte{
		// Record header
		0x16,       // type is 0x16 (handshake record)
		0x03, 0x04, // protocol version is "3.4" (also known as TLS 1.3)
	}

	var handshakeHeader []byte

	handshakeHeader = append(handshakeHeader, clientType[:]...)

	handshakeData := []byte{
		// Client Version
		0x03, 0x04, // protocol version is "3.4" (also known as TLS 1.3)
	}

	// Client random
	handshakeData = append(handshakeData, secret[:]...)

	handshakeData = append(handshakeData, []byte{
		// No Session ID
		0x00,

		// Cipher Suites
		0x00, 0x02,
	}...)
	handshakeData = append(handshakeData, tlsAes128GcmSha256[:]...)
	handshakeData = append(handshakeData, []byte{
		// Compression Methods
		0x01, // 1 byte of compression methods
		0x00, // no compression
	}...)

	encodedExtensions, err := clientExtensions(hostNames, keys)
	if err != nil {
		return nil, fmt.Errorf("failed encoding client hello message extensions: %w", err)
	}
	handshakeData = binary.BigEndian.AppendUint16(handshakeData, uint16(len(encodedExtensions)))
	handshakeData = append(handshakeData, encodedExtensions...)

	handshakeHeader, err = BigEndianAppend24(handshakeHeader, uint32(len(handshakeData)))
	if err != nil {
		return nil, fmt.Errorf("failed to append handshake header length: %w", err)
	}

	helloMessage = binary.BigEndian.AppendUint16(helloMessage, uint16(len(handshakeHeader)+len(handshakeData)))
	helloMessage = append(helloMessage, handshakeHeader...)
	helloMessage = append(helloMessage, handshakeData...)

	return helloMessage, nil
}
