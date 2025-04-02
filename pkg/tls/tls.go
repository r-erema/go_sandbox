package tls

import (
	"fmt"

	"golang.org/x/crypto/cryptobyte"
)

const (
	secretLength = 32
)

type (
	messageType               [1]byte
	cipherSuite               [2]byte
	supportedKeyExchangeGroup [2]byte
	signatureAlgorithm        [2]byte

	supportedTLSVersion [2]byte

	pskKeyExchangeMode [1]byte

	publicKey struct {
		payload       [32]byte
		exchangeGroup supportedKeyExchangeGroup
	}
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

var (
	clientType = messageType{0x01}
	serverType = messageType{0x02}

	tlsAes128GcmSha256 = cipherSuite{0x13, 0x01}

	x25519 = supportedKeyExchangeGroup{0x00, 0x1d}

	rsaPssRsaeSha512 = signatureAlgorithm{0x08, 0x06}

	tls13 = supportedTLSVersion{0x03, 0x04}

	pskWithECDHEKeyEstablishment = pskKeyExchangeMode{0x01}
)

func decodeRawHello(helloMessagePayload []byte) (*hello, error) {
	r := cryptobyte.String(helloMessagePayload)

	helloMsg := new(hello)

	r.ReadBytes(&helloMsg.recordType, 1)
	r.ReadBytes(&helloMsg.protocolVersion, 2)

	r.ReadUint16LengthPrefixed(&r)

	var buf []byte
	r.ReadBytes(&buf, 2)
	helloMsg.handshakeType = messageType(buf)

	r.ReadUint16LengthPrefixed(&r)

	r.ReadBytes(&helloMsg.version, 2)
	r.ReadBytes(&helloMsg.random, secretLength)

	r.Skip(1) // skip Session ID

	var ciphers cryptobyte.String
	switch helloMsg.handshakeType {
	case clientType:
		r.ReadUint16LengthPrefixed(&ciphers)
	case serverType:
		buf = []byte{}
		r.ReadBytes(&buf, 2)
		ciphers = buf
	}
	for !ciphers.Empty() {
		buf = []byte{}
		ciphers.ReadBytes(&buf, 2)
		helloMsg.cipherSuites = append(helloMsg.cipherSuites, cipherSuite(buf))
	}

	var compression cryptobyte.String
	r.ReadUint8LengthPrefixed(&compression)
	helloMsg.compression = compression

	r.ReadUint16LengthPrefixed(&r) // extensions length

	for !r.Empty() {
		var extType uint16
		r.ReadUint16(&extType)

		var extensionPayload cryptobyte.String
		r.ReadUint16LengthPrefixed(&extensionPayload)
		switch true {
		case extType == extensionTypeServerName && helloMsg.handshakeType == clientType:
			helloMsg.extensionServerName = decodeServerNameExtension(extensionPayload)
		case extType == extensionTypeEllipticCurvePointFormats && helloMsg.handshakeType == clientType:
			helloMsg.extensionECPointFormats = decodeECPointFormatsExtension(extensionPayload)
		case extType == extensionTypeSupportedGroups && helloMsg.handshakeType == clientType:
			helloMsg.extensionSupportedGroups = decodeSupportedGroupExtension(extensionPayload)
		case extType == extensionTypeSessionTicket && helloMsg.handshakeType == clientType:
			helloMsg.extensionSessionTicket = decodeSessionTicketExtension(extensionPayload)
		case extType == extensionTypeEncryptThenMAC && helloMsg.handshakeType == clientType:
			helloMsg.extensionEncryptThenMAC = decodeEncryptThenMACExtension(extensionPayload)
		case extType == extensionTypeExtendedMasterSecret && helloMsg.handshakeType == clientType:
			helloMsg.extensionExtendedMasterSecret = decodeExtendedMasterSecretExtension(extensionPayload)
		case extType == extensionTypeSignatureAlgorithms && helloMsg.handshakeType == clientType:
			helloMsg.extensionSignatureAlgorithms = decodeSignatureAlgorithmsExtension(extensionPayload)
		case extType == extensionTypePSKKeyExchangeModes && helloMsg.handshakeType == clientType:
			helloMsg.extensionPSKKeyExchangeModes = decodePSKKeyExchangeModesExtension(extensionPayload)
		case extType == extensionTypeSupportedTLSVersions:
			helloMsg.extensionSupportedTLSVersions = decodeSupportedTLSVersionsExtension(extensionPayload)
		case extType == extensionTypeKeyShare && helloMsg.handshakeType == clientType:
			helloMsg.extensionKeyShare = decodeClientKeyShareExtension(extensionPayload)
		case extType == extensionTypeKeyShare && helloMsg.handshakeType == serverType:
			helloMsg.extensionKeyShare = decodeServerKeyShareExtension(extensionPayload)
		default:
			return nil, NewError(fmt.Sprintf("unknown extension type: %d or wrong handshake type %q", extType, helloMsg.handshakeType))
		}
	}

	return helloMsg, nil
}
