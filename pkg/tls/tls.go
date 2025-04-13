package tls

import (
	"fmt"

	"golang.org/x/crypto/cryptobyte"
)

const (
	secretLength    = 32
	publicKeyLength = 32

	bytesCountForShortDataLength = 1
	bytesCountForLongDataLength  = 2
)

type hello struct {
	recordType,
	protocolVersion,

	version,
	random,
	compression []byte

	handshakeType messageType

	cipherSuites []cipherSuite

	extensionServerName
	extensionECPointFormats
	extensionSupportedGroups
	extensionSessionTicket
	extensionEncryptThenMAC
	extensionExtendedMasterSecret
	extensionSignatureAlgorithms
	extensionSupportedTLSVersions
	extensionPSKKeyExchangeModes
	extensionKeyShare
}

func decodeRawHello(helloMessagePayload []byte) (*hello, error) {
	raw := cryptobyte.String(helloMessagePayload)

	helloMsg := new(hello)

	raw.ReadBytes(&helloMsg.recordType, 1)
	raw.ReadBytes(&helloMsg.protocolVersion, bytesCountForLongDataLength)

	raw.ReadUint16LengthPrefixed(&raw)

	var buf []byte

	raw.ReadBytes(&buf, bytesCountForLongDataLength)
	helloMsg.handshakeType = messageType(buf)

	raw.ReadUint16LengthPrefixed(&raw)

	raw.ReadBytes(&helloMsg.version, bytesCountForLongDataLength)
	raw.ReadBytes(&helloMsg.random, secretLength)

	raw.Skip(1) // skip Session ID

	var ciphers cryptobyte.String

	switch helloMsg.handshakeType {
	case clientType():
		raw.ReadUint16LengthPrefixed(&ciphers)
	case serverType():
		buf = []byte{}
		raw.ReadBytes(&buf, bytesCountForLongDataLength)
		ciphers = buf
	}

	for !ciphers.Empty() {
		buf = []byte{}
		ciphers.ReadBytes(&buf, bytesCountForLongDataLength)
		helloMsg.cipherSuites = append(helloMsg.cipherSuites, cipherSuite(buf))
	}

	var compression cryptobyte.String

	raw.ReadUint8LengthPrefixed(&compression)
	helloMsg.compression = compression

	raw.ReadUint16LengthPrefixed(&raw) // extensions length

	if err := parseExtensions(raw, helloMsg); err != nil {
		return nil, fmt.Errorf("failed to parse extensions: %w", err)
	}

	return helloMsg, nil
}
