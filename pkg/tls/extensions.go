package tls

import (
	"encoding/binary"
	"fmt"

	"github.com/spf13/cast"
	"golang.org/x/crypto/cryptobyte"
)

const (
	extensionTypeServerName                = 0x00
	extensionTypeEllipticCurvePointFormats = 0x0b
	extensionTypeSupportedGroups           = 0x0a
	extensionTypeSessionTicket             = 0x23
	extensionTypeEncryptThenMAC            = 0x16
	extensionTypeExtendedMasterSecret      = 0x17
	extensionTypeSignatureAlgorithms       = 0x0d
	extensionTypeSupportedTLSVersions      = 0x2b
	extensionTypePSKKeyExchangeModes       = 0x2d
	extensionTypeKeyShare                  = 0x33
)

type (
	extensionServerName struct {
		entriesList []extensionServerNameListEntry
	}
	extensionServerNameListEntry struct {
		entryType uint8
		host      []byte
	}

	extensionECPointFormats struct {
		types [][1]byte
	}

	extensionSupportedGroups struct {
		groups []supportedKeyExchangeGroup
	}

	extensionSessionTicket struct {
		payload []byte
	}

	extensionEncryptThenMAC struct {
		payload []byte
	}

	extensionExtendedMasterSecret struct {
		payload []byte
	}

	extensionSignatureAlgorithms struct {
		algorithms []signatureAlgorithm
	}

	extensionSupportedTLSVersions struct {
		versions []supportedTLSVersion
	}

	extensionPSKKeyExchangeModes struct {
		modes []pskKeyExchangeMode
	}
	extensionKeyShare struct {
		publicKeys []publicKey
	}
)

func clientExtensions(hostNames []string, keys []publicKey) ([]byte, error) {
	var res []byte

	data, err := encodeServerNameExtension(hostNames)
	if err != nil {
		return nil, fmt.Errorf("failed to encode server name extension: %w", err)
	}

	res = append(res, data...)

	data, err = encodeECPointFormatsExtension()
	if err != nil {
		return nil, fmt.Errorf("failed to encode EC point formats extension: %w", err)
	}

	res = append(res, data...)

	data, err = encodeSupportedGroupsExtension([]supportedKeyExchangeGroup{
		x25519(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode supported groups extension: %w", err)
	}

	res = append(res, data...)
	res = append(res, encodeSessionTicketExtension()...)
	res = append(res, encodeEncryptThenMACExtension()...)
	res = append(res, encodeExtendedMasterSecretExtension()...)

	data, err = encodeSignatureAlgorithmsExtension([]signatureAlgorithm{rsaPssRsaeSha512()})
	if err != nil {
		return nil, fmt.Errorf("failed to encode signature algorithms: %w", err)
	}

	res = append(
		res,
		data...)

	data, err = encodeSupportedVersionsExtension([]supportedTLSVersion{tls13()})
	if err != nil {
		return nil, fmt.Errorf("failed to encode supported versions: %w", err)
	}

	res = append(res, data...)
	res = append(res, encodePSKKeyExchangeModesExtension()...)

	keyShareExtensionData, err := encodeClientKeyShareExtension(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to encode key share extension: %w", err)
	}

	res = append(res, keyShareExtensionData...)

	return res, nil
}

func serverExtensions(key publicKey) ([]byte, error) {
	var res []byte

	data, err := encodeSupportedVersionsExtension([]supportedTLSVersion{tls13()})
	if err != nil {
		return nil, fmt.Errorf("failed to encode supported versions extension: %w", err)
	}

	res = append(res, data...)

	data, err = encodeServerKeyShareExtension(key)
	if err != nil {
		return nil, fmt.Errorf("failed to encode key share extension: %w", err)
	}

	res = append(res, data...)

	return res, nil
}

func encodeServerNameExtension(hostNames []string) ([]byte, error) {
	const dnsHostNameType = 0

	var entries []byte

	extMetaLength, entryMetaLength, hostNameMetaLength := 4, 5, 3

	for _, host := range hostNames {
		hostNameEncoded := []byte(host)

		entry := make([]byte, entryMetaLength)

		hostLen, err := cast.ToUint16E(len(hostNameEncoded) + hostNameMetaLength)
		if err != nil {
			return nil, fmt.Errorf("failed to convert host length to uint16 type: %w", err)
		}

		hostNameLen, err := cast.ToUint16E(len(hostNameEncoded))
		if err != nil {
			return nil, fmt.Errorf("failed to convert host length to uint16 type: %w", err)
		}

		binary.BigEndian.PutUint16(entry[0:], hostLen)
		binary.BigEndian.PutUint16(entry[2:], dnsHostNameType)
		binary.BigEndian.PutUint16(entry[3:], hostNameLen)

		entry, err = binary.Append(entry, binary.BigEndian, hostNameEncoded)
		if err != nil {
			return nil, fmt.Errorf("appending hostname binary data error: %w", err)
		}

		entries, err = binary.Append(entries, binary.BigEndian, entry)
		if err != nil {
			return nil, fmt.Errorf("appending hostname entry error: %w", err)
		}
	}

	res := make([]byte, extMetaLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeServerName)

	entriesLen, err := cast.ToUint16E(len(entries))
	if err != nil {
		return nil, fmt.Errorf("failed to convert hastname entries length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], entriesLen)

	res, err = binary.Append(res, binary.BigEndian, entries)
	if err != nil {
		return nil, fmt.Errorf("appending hostname entries error: %w", err)
	}

	return res, nil
}

func decodeServerNameExtension(extensionPayload cryptobyte.String) extensionServerName {
	var res extensionServerName

	for !extensionPayload.Empty() {
		var listEntry extensionServerNameListEntry

		var listEntryRaw cryptobyte.String

		extensionPayload.ReadUint16LengthPrefixed(&listEntryRaw)

		listEntryRaw.ReadUint8(&listEntry.entryType)

		var listEntryPayload cryptobyte.String

		listEntryRaw.ReadUint16LengthPrefixed(&listEntryPayload)
		listEntry.host = listEntryPayload

		res.entriesList = append(res.entriesList, listEntry)
	}

	return res
}

func encodeECPointFormatsExtension() ([]byte, error) {
	const (
		uncompressed            = 0x0
		ansix962CompressedPrime = 0x1
		ansix962CompressedChar2 = 0x2
	)

	supportedFormatsLen := len(
		[]byte{uncompressed, ansix962CompressedPrime, ansix962CompressedChar2},
	)

	extMetaLength := 8

	res := make([]byte, extMetaLength)

	binary.BigEndian.PutUint16(res[0:], extensionTypeEllipticCurvePointFormats)

	formatsLen, err := cast.ToUint16E(supportedFormatsLen + bytesCountForShortDataLength)
	if err != nil {
		return nil, fmt.Errorf("failed to convert supported formats length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], formatsLen)

	res[4] = byte(supportedFormatsLen)
	res[5] = uncompressed
	res[6] = ansix962CompressedPrime
	res[7] = ansix962CompressedChar2

	return res, nil
}

func decodeECPointFormatsExtension(extensionPayload cryptobyte.String) extensionECPointFormats {
	var formats cryptobyte.String

	extensionPayload.ReadUint8LengthPrefixed(&formats)

	ext := extensionECPointFormats{
		types: make([][1]byte, len(formats)),
	}

	for i := range formats {
		ext.types[i] = [1]byte{formats[i]}
	}

	return ext
}

func encodeSupportedGroupsExtension(groups []supportedKeyExchangeGroup) ([]byte, error) {
	var entries []byte

	for _, group := range groups {
		entries = append(entries, group[:]...)
	}

	extMetaLength := 6
	res := make([]byte, extMetaLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSupportedGroups)

	entriesLen, err := cast.ToUint16E(len(entries))
	if err != nil {
		return nil, fmt.Errorf("failed to convert supported groups entries length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], entriesLen+bytesCountForLongDataLength)
	binary.BigEndian.PutUint16(res[4:], entriesLen)

	res, err = binary.Append(res, binary.BigEndian, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to convert supported groups entries error: %w", err)
	}

	return res, nil
}

func decodeSupportedGroupExtension(extensionPayload cryptobyte.String) extensionSupportedGroups {
	var groups cryptobyte.String

	extensionPayload.ReadUint16LengthPrefixed(&groups)

	ext := extensionSupportedGroups{}
	groupsMetaLength := 2

	for !groups.Empty() {
		var buf []byte

		groups.ReadBytes(&buf, groupsMetaLength)
		ext.groups = append(ext.groups, supportedKeyExchangeGroup(buf))
	}

	return ext
}

func decodeSessionTicketExtension(extensionPayload cryptobyte.String) extensionSessionTicket {
	return extensionSessionTicket{payload: extensionPayload}
}

func decodeEncryptThenMACExtension(extensionPayload cryptobyte.String) extensionEncryptThenMAC {
	return extensionEncryptThenMAC{payload: extensionPayload}
}

func decodeExtendedMasterSecretExtension(
	extensionPayload cryptobyte.String,
) extensionExtendedMasterSecret {
	return extensionExtendedMasterSecret{payload: extensionPayload}
}

func decodeSignatureAlgorithmsExtension(
	extensionPayload cryptobyte.String,
) extensionSignatureAlgorithms {
	var algorithms cryptobyte.String

	extensionPayload.ReadUint16LengthPrefixed(&algorithms)

	ext := extensionSignatureAlgorithms{}

	algorithmsMetaLength := 2

	for !algorithms.Empty() {
		var buf []byte

		algorithms.ReadBytes(&buf, algorithmsMetaLength)
		ext.algorithms = append(ext.algorithms, signatureAlgorithm(buf))
	}

	return ext
}

func decodeSupportedTLSVersionsExtension(
	extensionPayload cryptobyte.String,
) extensionSupportedTLSVersions {
	var versions cryptobyte.String

	extensionPayload.ReadUint8LengthPrefixed(&versions)

	ext := extensionSupportedTLSVersions{}

	for !versions.Empty() {
		var buf []byte

		versions.ReadBytes(&buf, bytesCountForLongDataLength)
		ext.versions = append(ext.versions, supportedTLSVersion(buf))
	}

	return ext
}

func decodePSKKeyExchangeModesExtension(
	extensionPayload cryptobyte.String,
) extensionPSKKeyExchangeModes {
	var versions cryptobyte.String

	extensionPayload.ReadUint8LengthPrefixed(&versions)

	ext := extensionPSKKeyExchangeModes{}

	for !versions.Empty() {
		var buf []byte

		versions.ReadBytes(&buf, 1)
		ext.modes = append(ext.modes, pskKeyExchangeMode(buf))
	}

	return ext
}

func decodeClientKeyShareExtension(extensionPayload cryptobyte.String) extensionKeyShare {
	var keys cryptobyte.String

	extensionPayload.ReadUint16LengthPrefixed(&keys)

	ext := extensionKeyShare{}

	for !keys.Empty() {
		var buf []byte

		var key publicKey

		keys.ReadBytes(&buf, bytesCountForLongDataLength)
		key.exchangeGroup = supportedKeyExchangeGroup(buf)

		keys.ReadUint16LengthPrefixed(&keys)
		keys.ReadBytes(&buf, keyLength)
		key.payload = [keyLength]byte(buf)

		ext.publicKeys = append(ext.publicKeys, key)
	}

	return ext
}

func decodeServerKeyShareExtension(extensionPayload cryptobyte.String) extensionKeyShare {
	ext := extensionKeyShare{}

	var buf []byte

	key := publicKey{}

	extensionPayload.ReadBytes(&buf, bytesCountForLongDataLength)
	key.exchangeGroup = supportedKeyExchangeGroup(buf)

	extensionPayload.ReadUint16LengthPrefixed(&extensionPayload)
	key.payload = [keyLength]byte(extensionPayload)

	ext.publicKeys = append(ext.publicKeys, key)

	return ext
}

func encodeSessionTicketExtension() []byte {
	extLength := 4
	res := make([]byte, extLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSessionTicket)

	return res
}

func encodeEncryptThenMACExtension() []byte {
	return []byte{
		0x0, extensionTypeEncryptThenMAC,
		0x0, 0x0,
	}
}

func encodeExtendedMasterSecretExtension() []byte {
	return []byte{
		0x0, extensionTypeExtendedMasterSecret,
		0x0, 0x0,
	}
}

func encodeSignatureAlgorithmsExtension(algos []signatureAlgorithm) ([]byte, error) {
	var entries []byte

	for _, algoBytes := range algos {
		entries = append(entries, algoBytes[:]...)
	}

	initialLength := 6

	res := make([]byte, initialLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSignatureAlgorithms)

	entriesLen, err := cast.ToUint16E(len(entries))
	if err != nil {
		return nil, fmt.Errorf("failed to convert signature algorithms entries length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], entriesLen+bytesCountForLongDataLength)
	binary.BigEndian.PutUint16(res[4:], entriesLen)

	res, err = binary.Append(res, binary.BigEndian, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to append entries data: %w", err)
	}

	return res, nil
}

func encodeSupportedVersionsExtension(versions []supportedTLSVersion) ([]byte, error) {
	var entries []byte

	for _, versionBytes := range versions {
		entries = append(entries, versionBytes[:]...)
	}

	initialLength := 5
	res := make([]byte, initialLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSupportedTLSVersions)

	entriesLen, err := cast.ToUint16E(len(entries))
	if err != nil {
		return nil, fmt.Errorf("failed to convert supported versions entries length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], entriesLen+bytesCountForShortDataLength)
	res[4] = byte(len(entries))

	res, err = binary.Append(res, binary.BigEndian, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to append entries data: %w", err)
	}

	return res, nil
}

func encodePSKKeyExchangeModesExtension() []byte {
	ext := []byte{
		0x00, extensionTypePSKKeyExchangeModes,
		0x0, 0x02,
		0x01,
	}

	exchangeMode := pskWithECDHEKeyEstablishment()
	ext = append(ext, exchangeMode[:]...)

	return ext
}

func encodeClientKeyShareExtension(publicKeys []publicKey) ([]byte, error) {
	var entries []byte

	for _, publicKey := range publicKeys {
		keyPayload := make([]byte, bytesCountForLongDataLength)

		keyLen, err := cast.ToUint16E(len(publicKey.payload))
		if err != nil {
			return nil, fmt.Errorf("failed to convert public key length to uint16 type: %w", err)
		}

		binary.BigEndian.PutUint16(keyPayload[0:], keyLen)

		keyPayload, err = binary.Append(keyPayload, binary.BigEndian, publicKey.payload[:])
		if err != nil {
			return nil, fmt.Errorf("failed to append entries data: %w", err)
		}

		entry := make([]byte, bytesCountForLongDataLength)

		entryPayload, err := binary.Append(publicKey.exchangeGroup[:], binary.BigEndian, keyPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to append bytes to the entry payload: %w", err)
		}

		entryLen, err := cast.ToUint16E(len(entryPayload))
		if err != nil {
			return nil, fmt.Errorf("failed to convert public key entry payload length to uint16 type: %w", err)
		}

		binary.BigEndian.PutUint16(entry[0:], entryLen)

		entry, err = binary.Append(entry, binary.BigEndian, entryPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to append entries data: %w", err)
		}

		entries, err = binary.Append(entries, binary.BigEndian, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to append entries data: %w", err)
		}
	}

	initialLength := 4
	res := make([]byte, initialLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeKeyShare)

	entriesLen, err := cast.ToUint16E(len(entries))
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key entries length to uint16 type: %w", err)
	}

	binary.BigEndian.PutUint16(res[2:], entriesLen)

	res, err = binary.Append(res, binary.BigEndian, entries)
	if err != nil {
		return nil, fmt.Errorf("failed to append entries data: %w", err)
	}

	return res, nil
}

func encodeServerKeyShareExtension(publicKey publicKey) ([]byte, error) {
	initialLength := 8

	res := make([]byte, initialLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeKeyShare)
	binary.BigEndian.PutUint16(
		res[2:],
		uint16(
			len(
				publicKey.exchangeGroup,
			)+len(
				publicKey.payload,
			)+bytesCountForLongDataLength,
		),
	)
	binary.BigEndian.PutUint16(res[4:], binary.BigEndian.Uint16(publicKey.exchangeGroup[:]))
	binary.BigEndian.PutUint16(res[6:], uint16(len(publicKey.payload)))

	res, err := binary.Append(res, binary.BigEndian, publicKey.payload[:])
	if err != nil {
		return nil, fmt.Errorf("failed to append public key payload data: %w", err)
	}

	return res, nil
}

func parseExtensions(raw cryptobyte.String, helloMsg *hello) error {
	clientDecoders := map[uint16]func(*hello, cryptobyte.String){
		extensionTypeServerName: func(helloMsg *hello, payload cryptobyte.String) {
			helloMsg.extensionServerName = decodeServerNameExtension(payload)
		},
		extensionTypeEllipticCurvePointFormats: func(helloMsg *hello, payload cryptobyte.String) {
			helloMsg.extensionECPointFormats = decodeECPointFormatsExtension(payload)
		},
		extensionTypeSupportedGroups: func(helloMsg *hello, payload cryptobyte.String) {
			helloMsg.extensionSupportedGroups = decodeSupportedGroupExtension(payload)
		},
		extensionTypeSessionTicket: func(h *hello, payload cryptobyte.String) {
			h.extensionSessionTicket = decodeSessionTicketExtension(payload)
		},
		extensionTypeEncryptThenMAC: func(h *hello, payload cryptobyte.String) {
			h.extensionEncryptThenMAC = decodeEncryptThenMACExtension(payload)
		},
		extensionTypeExtendedMasterSecret: func(h *hello, payload cryptobyte.String) {
			h.extensionExtendedMasterSecret = decodeExtendedMasterSecretExtension(payload)
		},
		extensionTypeSignatureAlgorithms: func(h *hello, payload cryptobyte.String) {
			h.extensionSignatureAlgorithms = decodeSignatureAlgorithmsExtension(payload)
		},
		extensionTypePSKKeyExchangeModes: func(h *hello, payload cryptobyte.String) {
			h.extensionPSKKeyExchangeModes = decodePSKKeyExchangeModesExtension(payload)
		},
		extensionTypeSupportedTLSVersions: func(h *hello, payload cryptobyte.String) {
			h.extensionSupportedTLSVersions = decodeSupportedTLSVersionsExtension(payload)
		},
		extensionTypeKeyShare: func(h *hello, payload cryptobyte.String) {
			h.extensionKeyShare = decodeClientKeyShareExtension(payload)
		},
	}

	serverDecoders := map[uint16]func(*hello, cryptobyte.String){
		extensionTypeSupportedTLSVersions: func(helloMsg *hello, payload cryptobyte.String) {
			helloMsg.extensionSupportedTLSVersions = decodeSupportedTLSVersionsExtension(payload)
		},
		extensionTypeKeyShare: func(helloMsg *hello, payload cryptobyte.String) {
			helloMsg.extensionKeyShare = decodeServerKeyShareExtension(payload)
		},
	}

	var decoders map[uint16]func(*hello, cryptobyte.String)

	switch {
	case helloMsg.handshakeType == clientType():
		decoders = clientDecoders
	case helloMsg.handshakeType == serverType():
		decoders = serverDecoders
	default:
		return NewError(fmt.Sprintf("unknown handshake type: %d", helloMsg.handshakeType))
	}

	var (
		extType          uint16
		extensionPayload cryptobyte.String
	)

	for !raw.Empty() {
		raw.ReadUint16(&extType)
		raw.ReadUint16LengthPrefixed(&extensionPayload)

		if decoder, ok := decoders[extType]; ok {
			decoder(helloMsg, extensionPayload)

			continue
		}

		return NewError(fmt.Sprintf("unknown extension type: %d", extType))
	}

	return nil
}
