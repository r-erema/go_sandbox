package tls

import (
	"encoding/binary"
	"fmt"

	"github.com/spf13/cast"
)

func encodeServerHello(key publicKey, secret [32]byte, clientCipher cipherSuite) ([]byte, error) {
	helloMessage := []byte{
		// Record header
		0x16,       // type is 0x16 (handshake record)
		0x03, 0x03, // protocol version is "3.4" (also known as TLS 1.3)
	}

	var handshakeHeader []byte

	helloServerType := serverType()
	handshakeHeader = append(handshakeHeader, helloServerType[:]...)

	handshakeData := []byte{
		// Server Version
		0x03, 0x03, // protocol version is "3.4" (also known as TLS 1.3)
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

	encodedExtensions, err := serverExtensions(key)
	if err != nil {
		return nil, fmt.Errorf("failed to encode server extensions: %w", err)
	}

	extensionsLen, err := cast.ToUint16E(len(encodedExtensions))
	if err != nil {
		return nil, fmt.Errorf("failed to convert extensions length to uint16 type: %w", err)
	}

	handshakeData = binary.BigEndian.AppendUint16(handshakeData, extensionsLen)
	handshakeData = append(handshakeData, encodedExtensions...)

	handshakeDataLen, err := cast.ToUint32E(len(handshakeData))
	if err != nil {
		return nil, fmt.Errorf("failed to convert handshake data length to uint32 type: %w", err)
	}

	handshakeHeader, err = BigEndianAppend24(handshakeHeader, handshakeDataLen)
	if err != nil {
		return nil, fmt.Errorf("failed to append handshake header length: %w", err)
	}

	handshakeLen, err := cast.ToUint16E(len(handshakeHeader) + len(handshakeData))
	if err != nil {
		return nil, fmt.Errorf("failed to convert handshake length to uint16 type: %w", err)
	}

	helloMessage = binary.BigEndian.AppendUint16(
		helloMessage,
		handshakeLen,
	)
	helloMessage = append(helloMessage, handshakeHeader...)
	helloMessage = append(helloMessage, handshakeData...)

	return helloMessage, nil
}
