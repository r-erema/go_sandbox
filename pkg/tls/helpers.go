package tls

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"fmt"

	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"github.com/spf13/cast"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

func Rand32Bytes() ([secretLength]byte, error) {
	buf := make([]byte, secretLength)
	if _, err := rand.Read(buf); err != nil {
		return [secretLength]byte{}, fmt.Errorf("rand read error: %w", err)
	}

	return [secretLength]byte(buf), nil
}

func GeneratePrivateAndPublicKeys() ([secretLength]byte, [secretLength]byte, error) {
	privateKey, err := Rand32Bytes()
	if err != nil {
		return [secretLength]byte{}, [secretLength]byte{}, fmt.Errorf("failed generating private key: %w", err)
	}

	pubKey, err := curve25519.X25519(privateKey[:], curve25519.Basepoint)
	if err != nil {
		return [secretLength]byte{}, [secretLength]byte{}, fmt.Errorf("failed generating public key: %w", err)
	}

	return privateKey, [secretLength]byte(pubKey), nil
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

func read(socketFD int) ([]byte, error) {
	header, err := syscall.Read(socketFD, headerLength)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	length := binary.BigEndian.Uint16(header[payloadStartIndexInHeader:])

	payload, err := syscall.Read(socketFD, int(length))
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	res, err := binary.Append(header, binary.BigEndian, payload)
	if err != nil {
		return nil, fmt.Errorf("failed concat header and payload: %w", err)
	}

	return res, nil
}

func WriteAppData(socketFD int, message, clientAppKey ClientApplicationKey, initVector []byte) error {
	message = append(message, messageTypeApplication)

	messageLen, err := cast.ToUint16E(len(message) + authTagLength)
	if err != nil {
		return fmt.Errorf("failed to convert message length to uint16 type: %w", err)
	}

	additional := binary.BigEndian.AppendUint16([]byte{messageTypeApplication, 0x03, 0x03}, messageLen)

	encryptedData, err := encrypt(clientAppKey, initVector, message, additional)
	if err != nil {
		return fmt.Errorf("failed encrypting client application key: %w", err)
	}

	if err = syscall.Write(socketFD, encryptedData); err != nil {
		return fmt.Errorf("failed sending encrypted app data finished message: %w", err)
	}

	return nil
}

func ReadAppData(socketFD int, serverAppKey ServerApplicationKey, initVector []byte) ([]byte, error) {
	buf, err := read(socketFD)
	if err != nil {
		return nil, fmt.Errorf("failed reading server handshake message: %w", err)
	}

	buf, err = Decrypt(serverAppKey, initVector, buf)
	if err != nil {
		return nil, fmt.Errorf("failed decrypting server handshake message: %w", err)
	}

	return buf, nil
}

func deriveSecret(secret []byte, label string, transcriptMessages []byte) ([]byte, error) {
	hash := sha512.Sum384(transcriptMessages)

	hashLen, err := cast.ToUint16E(len(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to convert hash length to uint16 type: %w", err)
	}

	res, err := hkdfExpandLabel(secret, label, hash[:], hashLen)
	if err != nil {
		return nil, fmt.Errorf("failed deriving secret: %w", err)
	}

	return res, nil
}

func hkdfExpandLabel(secret []byte, label string, transcriptHash []byte, hashLen uint16) ([]byte, error) {
	var hkdfLabel cryptobyte.Builder

	hkdfLabel.AddUint16(hashLen)
	hkdfLabel.AddUint8LengthPrefixed(func(child *cryptobyte.Builder) {
		child.AddBytes([]byte("tls13 "))
		child.AddBytes([]byte(label))
	})
	hkdfLabel.AddUint8LengthPrefixed(func(child *cryptobyte.Builder) {
		child.AddBytes(transcriptHash)
	})

	buf, err := hkdfLabel.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to build bytes: %w", err)
	}

	reader := hkdf.Expand(sha512.New384, secret, buf)

	buf = make([]byte, hashLen)

	_, err = reader.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read expanded label bytes: %w", err)
	}

	return buf, nil
}

func Decrypt(key, initVector, wrapper []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	additional, ciphertext := wrapper[:5], wrapper[5:]

	plain, err := gcm.Open(nil, initVector, ciphertext, additional)
	if err != nil {
		return nil, fmt.Errorf("failed to Decrypt ciphertext: %w", err)
	}

	return plain, nil
}

func IV(counter uint8, iv []byte) []byte {
	res := make([]byte, len(iv))
	copy(res, iv)

	offset := 12

	for i := range offset {
		res[len(res)-i-1] ^= counter >> cast.ToUint(offset*i)
	}

	return res
}

func encrypt(key, initVector, plaintext, additional []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, initVector, plaintext, additional)

	res, err := binary.Append(additional, binary.BigEndian, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to append ciphertext and additional: %w", err)
	}

	return res, nil
}
