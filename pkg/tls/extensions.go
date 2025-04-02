package tls

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/cryptobyte"
)

const (
	bytesCountForExtensionShortDataLength = 1
	bytesCountForExtensionLongDataLength  = 2

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

	serverNameExtData, err := encodeServerNameExtension(hostNames)
	if err != nil {
		return nil, fmt.Errorf("failed to encode server name extension: %w", err)
	}

	res = append(res, serverNameExtData...)
	res = append(res, encodeECPointFormatsExtension()...)
	res = append(res, encodeSupportedGroupsExtension([]supportedKeyExchangeGroup{
		x25519,
	})...)
	res = append(res, encodeSessionTicketExtension()...)
	res = append(res, encodeEncryptThenMACExtension()...)
	res = append(res, encodeExtendedMasterSecretExtension()...)
	res = append(res, encodeSignatureAlgorithmsExtension([]signatureAlgorithm{rsaPssRsaeSha512})...)
	res = append(res, encodeSupportedVersionsExtension([]supportedTLSVersion{tls13})...)
	res = append(res, encodePSKKeyExchangeModesExtension()...)
	res = append(res, encodeClientKeyShareExtension(keys)...)

	return res, nil
}

func serverExtensions(key publicKey) []byte {
	var res []byte
	res = append(res, encodeSupportedVersionsExtension([]supportedTLSVersion{tls13})...)
	res = append(res, encodeServerKeyShareExtension(key)...)

	return res
}

func encodeServerNameExtension(hostNames []string) ([]byte, error) {
	const dnsHostNameType = 0

	var entries []byte

	extMetaLength, entryMetaLength, hostNameMetaLength := 4, 5, 3

	for _, host := range hostNames {
		hostNameEncoded := []byte(host)

		hostNameLen := len(hostNameEncoded)

		entry := make([]byte, entryMetaLength)
		binary.BigEndian.PutUint16(entry[0:], uint16(hostNameLen+hostNameMetaLength))
		binary.BigEndian.PutUint16(entry[2:], dnsHostNameType)
		binary.BigEndian.PutUint16(entry[3:], uint16(hostNameLen))

		entry, err := binary.Append(entry, binary.BigEndian, hostNameEncoded)
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
	binary.BigEndian.PutUint16(res[2:], uint16(len(entries)))
	res, err := binary.Append(res, binary.BigEndian, entries)
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

func encodeECPointFormatsExtension() []byte {
	const (
		uncompressed            = 0x0
		ansix962CompressedPrime = 0x1
		ansix962CompressedChar2 = 0x2
	)

	supportedFormatsLen := len([]byte{uncompressed, ansix962CompressedPrime, ansix962CompressedChar2})

	extMetaLength := 8

	res := make([]byte, extMetaLength)

	binary.BigEndian.PutUint16(res[0:], extensionTypeEllipticCurvePointFormats)
	binary.BigEndian.PutUint16(res[2:], uint16(supportedFormatsLen+bytesCountForExtensionShortDataLength))
	res[4] = byte(supportedFormatsLen)
	res[5] = uncompressed
	res[6] = ansix962CompressedPrime
	res[7] = ansix962CompressedChar2

	return res
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

func encodeSupportedGroupsExtension(groups []supportedKeyExchangeGroup) []byte {
	var entries []byte

	for _, group := range groups {
		entries = append(entries, group[:]...)
	}

	extMetaLength := 6
	res := make([]byte, extMetaLength)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSupportedGroups)
	binary.BigEndian.PutUint16(res[2:], uint16(len(entries)+bytesCountForExtensionLongDataLength))
	binary.BigEndian.PutUint16(res[4:], uint16(len(entries)))
	res = append(res, entries...)

	return res
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

func decodeExtendedMasterSecretExtension(extensionPayload cryptobyte.String) extensionExtendedMasterSecret {
	return extensionExtendedMasterSecret{payload: extensionPayload}
}

func decodeSignatureAlgorithmsExtension(extensionPayload cryptobyte.String) extensionSignatureAlgorithms {
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

func decodeSupportedTLSVersionsExtension(extensionPayload cryptobyte.String) extensionSupportedTLSVersions {
	var versions cryptobyte.String
	extensionPayload.ReadUint8LengthPrefixed(&versions)

	ext := extensionSupportedTLSVersions{}

	for !versions.Empty() {
		var buf []byte
		versions.ReadBytes(&buf, 2)
		ext.versions = append(ext.versions, supportedTLSVersion(buf))
	}

	return ext
}

func decodePSKKeyExchangeModesExtension(extensionPayload cryptobyte.String) extensionPSKKeyExchangeModes {
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

		keys.ReadBytes(&buf, 2)
		key.exchangeGroup = supportedKeyExchangeGroup(buf)
		keys.ReadUint16LengthPrefixed(&keys)
		keys.ReadBytes(&buf, 32)
		key.payload = [32]byte(buf)

		ext.publicKeys = append(ext.publicKeys, key)
	}

	return ext
}

func decodeServerKeyShareExtension(extensionPayload cryptobyte.String) extensionKeyShare {
	ext := extensionKeyShare{}

	var buf []byte
	key := publicKey{}
	extensionPayload.ReadBytes(&buf, 2)
	key.exchangeGroup = supportedKeyExchangeGroup(buf)

	extensionPayload.ReadUint16LengthPrefixed(&extensionPayload)
	key.payload = [32]byte(extensionPayload)

	ext.publicKeys = append(ext.publicKeys, key)

	return ext
}

func encodeSessionTicketExtension() []byte {
	res := make([]byte, 4)
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

func encodeSignatureAlgorithmsExtension(algos []signatureAlgorithm) []byte {
	var entries []byte

	for _, algoBytes := range algos {
		entries = append(entries, algoBytes[:]...)
	}

	res := make([]byte, 6)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSignatureAlgorithms)
	binary.BigEndian.PutUint16(res[2:], uint16(len(entries)+bytesCountForExtensionLongDataLength))
	binary.BigEndian.PutUint16(res[4:], uint16(len(entries)))
	res = append(res, entries...)

	return res
}

func encodeSupportedVersionsExtension(versions []supportedTLSVersion) []byte {
	var entries []byte

	for _, versionBytes := range versions {
		entries = append(entries, versionBytes[:]...)
	}

	res := make([]byte, 5)
	binary.BigEndian.PutUint16(res[0:], extensionTypeSupportedTLSVersions)
	binary.BigEndian.PutUint16(res[2:], uint16(len(entries)+bytesCountForExtensionShortDataLength))
	res[4] = byte(len(entries))
	res = append(res, entries...)

	return res
}

func encodePSKKeyExchangeModesExtension() []byte {
	ext := []byte{
		0x00, extensionTypePSKKeyExchangeModes,
		0x0, 0x02,
		0x01,
	}
	ext = append(ext, pskWithECDHEKeyEstablishment[:]...)

	return ext
}

func encodeClientKeyShareExtension(publicKeys []publicKey) []byte {
	var entries []byte

	for _, publicKey := range publicKeys {
		keyPayload := make([]byte, 2)
		binary.BigEndian.PutUint16(keyPayload[0:], uint16(len(publicKey.payload)))
		keyPayload = append(keyPayload, publicKey.payload[:]...)

		entry := make([]byte, 2)
		entryPayload := append(publicKey.exchangeGroup[:], keyPayload...)
		binary.BigEndian.PutUint16(entry[0:], uint16(len(entryPayload)))
		entry = append(entry, entryPayload...)

		entries = append(entries, entry...)
	}

	res := make([]byte, 4)
	binary.BigEndian.PutUint16(res[0:], extensionTypeKeyShare)
	binary.BigEndian.PutUint16(res[2:], uint16(len(entries)))
	res = append(res, entries...)

	return res
}

func encodeServerKeyShareExtension(publicKey publicKey) []byte {
	res := make([]byte, 8)
	binary.BigEndian.PutUint16(res[0:], extensionTypeKeyShare)
	binary.BigEndian.PutUint16(res[2:], uint16(len(publicKey.exchangeGroup)+len(publicKey.payload)+bytesCountForExtensionLongDataLength))
	binary.BigEndian.PutUint16(res[4:], binary.BigEndian.Uint16(publicKey.exchangeGroup[:]))
	binary.BigEndian.PutUint16(res[6:], uint16(len(publicKey.payload)))
	res = append(res, publicKey.payload[:]...)

	return res
}
