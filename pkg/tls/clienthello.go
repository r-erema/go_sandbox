package tls

import (
	"encoding/binary"
	"fmt"

	"github.com/spf13/cast"
)

func encodeClientHello(hostNames []string, keys []publicKey, secret [32]byte) ([]byte, error) {
	helloMessage := []byte{
		// Record header
		0x16,       // type is 0x16 (handshake record)
		0x03, 0x03, // protocol version is "3.4" (also known as TLS 1.3)
	}

	var handshakeHeader []byte

	helloClientType := clientType()
	handshakeHeader = append(handshakeHeader, helloClientType[:]...)

	handshakeData := []byte{
		// Client Version
		0x03, 0x03, // protocol version is "3.4" (also known as TLS 1.3)
	}

	// Client random
	handshakeData = append(handshakeData, secret[:]...)

	handshakeData = append(handshakeData, []byte{
		// No Session ID
		0x00,

		// Cipher Suites
		0x00, 0x02,
	}...)

	cipher := tlsAes256GcmSha384()
	// cipher := tlsAes128GcmSha256()
	handshakeData = append(handshakeData, cipher[:]...)
	handshakeData = append(handshakeData, []byte{
		// Compression Methods
		0x01, // 1 byte of compression methods
		0x00, // no compression
	}...)

	encodedExtensions, err := clientExtensions(hostNames, keys)
	if err != nil {
		return nil, fmt.Errorf("failed encoding client hello message extensions: %w", err)
	}

	encodedExtensionsLen, err := cast.ToUint16E(len(encodedExtensions))
	if err != nil {
		return nil, fmt.Errorf("failed to convert encoded extensions length to uint16 type: %w", err)
	}

	handshakeData = binary.BigEndian.AppendUint16(handshakeData, encodedExtensionsLen)
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

	helloMessage = binary.BigEndian.AppendUint16(helloMessage, handshakeLen)
	helloMessage = append(helloMessage, handshakeHeader...)
	helloMessage = append(helloMessage, handshakeData...)

	return helloMessage, nil
}
