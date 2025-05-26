package tls

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/r-erema/go_sendbox/utils/os/syscall"
	"github.com/spf13/cast"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

func HandshakeClientSide(socketFD int, hosts []string, privateKey, clientPublicKey [32]byte) (
	ClientApplicationKey, ClientApplicationIV, ServerApplicationKey, ServerApplicationIV, error,
) {
	clientHelloRaw, serverHelloRaw, err := helloRaws(socketFD, hosts, clientPublicKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get client and server hello raws: %w", err)
	}

	serverHello, err := decodeRawHello(serverHelloRaw)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed decoding server hello message: %w", err)
	}

	handshakeSecret, err := deriveHandshakeSecret(privateKey, serverHello.extensionKeyShare.publicKeys[0].payload)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving handshake secret: %w", err)
	}

	handshakeMessages := slices.Concat(clientHelloRaw[headerLength:], serverHelloRaw[headerLength:])

	clientHandshakeSecret, err := deriveSecret(handshakeSecret, "c hs traffic", handshakeMessages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client handshake secret: %w", err)
	}

	clientHandshakeKey, clientHandshakeIV, serverHandshakeKey, serverHandshakeIV, err := handshakeKeys(
		clientHandshakeSecret,
		handshakeMessages,
		handshakeSecret,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed getting handshake keys: %w", err)
	}

	serverEncryptedExtension, serverCert, serverCertVerify, serverHandshake, err := serverData(
		socketFD,
		serverHandshakeKey,
		serverHandshakeIV,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed getting server data: %w", err)
	}

	handshakeMessages = slices.Concat(
		clientHelloRaw[headerLength:],
		serverHelloRaw[headerLength:],
		serverEncryptedExtension[:len(serverEncryptedExtension)-1],
		serverCert[:len(serverCert)-1],
		serverCertVerify[:len(serverCertVerify)-1],
		serverHandshake[:len(serverHandshake)-1],
	)

	if err = sendChangeCipherSuite(socketFD); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed sending cipher suite: %w", err)
	}

	msg, err := clientHandshakeFinishedMessage(
		clientHandshakeSecret,
		handshakeMessages,
		clientHandshakeKey,
		clientHandshakeIV,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed create handshake finished message: %w", err)
	}

	if err = syscall.Write(socketFD, msg); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed sending handshake finished message: %w", err)
	}

	return applicationKeys(handshakeSecret, handshakeMessages)
}

func helloRaws(socketFD int, hosts []string, clientPublicKey [32]byte) ([]byte, []byte, error) {
	secret, err := Rand32Bytes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed generating secret: %w", err)
	}

	clientHelloRaw, err := encodeClientHello(hosts, []publicKey{{
		payload:       clientPublicKey,
		exchangeGroup: x25519(),
	}}, secret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed encoding client hello message: %w", err)
	}

	if err = syscall.Write(socketFD, clientHelloRaw); err != nil {
		return nil, nil, fmt.Errorf("failed to write to file: %w", err)
	}

	serverHelloRaw, err := read(socketFD)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read connection of socket fd `%d`: %w", socketFD, err)
	}

	return clientHelloRaw, serverHelloRaw, nil
}

func deriveHandshakeSecret(clientPrivateKey, serverPublicKey [32]byte) ([]byte, error) {
	sharedSecret, err := curve25519.X25519(clientPrivateKey[:], serverPublicKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed generating shared secret: %w", err)
	}

	earlySecret := hkdf.Extract(sha512.New384, make([]byte, sha384Length), make([]byte, sha384Length))

	derivedEarlySecret, err := deriveSecret(earlySecret, "derived", []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed deriving secret: %w", err)
	}

	return hkdf.Extract(sha512.New384, sharedSecret, derivedEarlySecret), nil
}

func serverData(socketFD int, serverHandshakeKey, serverHandshakeIV []byte) ([]byte, []byte, []byte, []byte, error) {
	if err := readCipherSpecFromServer(socketFD); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed reading ciper spec from server: %w", err)
	}

	serverEncryptedExtension, err := read(socketFD)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed reading server encrypted extension: %w", err)
	}

	callCounter := uint8(0)

	serverEncryptedExtension, err = Decrypt(
		serverHandshakeKey,
		IV(callCounter, serverHandshakeIV),
		serverEncryptedExtension,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed decrypting server encrypted extension: %w", err)
	}

	callCounter++

	serverCert, err := read(socketFD)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed reading server certificate: %w", err)
	}

	serverCert, err = Decrypt(serverHandshakeKey, IV(callCounter, serverHandshakeIV), serverCert)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed decrypting server certificate: %w", err)
	}

	callCounter++

	serverCertVerify, err := read(socketFD)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed reading server certificate verify: %w", err)
	}

	serverCertVerify, err = Decrypt(serverHandshakeKey, IV(callCounter, serverHandshakeIV), serverCertVerify)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed decrypting server certificate verify: %w", err)
	}

	callCounter++

	serverHandshake, err := read(socketFD)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed reading server handshake message: %w", err)
	}

	serverHandshake, err = Decrypt(serverHandshakeKey, IV(callCounter, serverHandshakeIV), serverHandshake)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed decrypting server handshake message: %w", err)
	}

	return serverEncryptedExtension, serverCert, serverCertVerify, serverHandshake, nil
}

func handshakeKeys(
	clientHandshakeSecret, handshakeMessages, handshakeSecret []byte,
) ([]byte, []byte, []byte, []byte, error) {
	clientHandshakeKey, err := hkdfExpandLabel(clientHandshakeSecret, "key", []byte{}, keyLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client handshake key: %w", err)
	}

	clientHandshakeIV, err := hkdfExpandLabel(clientHandshakeSecret, "iv", []byte{}, ivLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client handshake IV: %w", err)
	}

	serverHandshakeSecret, err := deriveSecret(handshakeSecret, "s hs traffic", handshakeMessages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server handshake secret: %w", err)
	}

	serverHandshakeKey, err := hkdfExpandLabel(serverHandshakeSecret, "key", []byte{}, keyLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server handshake key: %w", err)
	}

	serverHandshakeIV, err := hkdfExpandLabel(serverHandshakeSecret, "iv", []byte{}, ivLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server handshake IV: %w", err)
	}

	return clientHandshakeKey, clientHandshakeIV, serverHandshakeKey, serverHandshakeIV, nil
}

func applicationKeys(
	handshakeSecret, handshakeMessages []byte,
) (ClientApplicationKey, ClientApplicationIV, ServerApplicationKey, ServerApplicationIV, error) {
	derivedSecret, err := deriveSecret(handshakeSecret, "derived", []byte{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving derived secret: %w", err)
	}

	masterSecret := hkdf.Extract(sha512.New384, make([]byte, sha384Length), derivedSecret)

	clientApplicationSecret, err := deriveSecret(masterSecret, "c ap traffic", handshakeMessages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client application secret: %w", err)
	}

	clientApplicationKey, err := hkdfExpandLabel(clientApplicationSecret, "key", []byte{}, keyLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client application key: %w", err)
	}

	clientApplicationIV, err := hkdfExpandLabel(clientApplicationSecret, "iv", []byte{}, ivLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving client application IV: %w", err)
	}

	serverApplicationSecret, err := deriveSecret(masterSecret, "s ap traffic", handshakeMessages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server application secret: %w", err)
	}

	serverApplicationKey, err := hkdfExpandLabel(serverApplicationSecret, "key", []byte{}, keyLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server application key: %w", err)
	}

	serverApplicationIV, err := hkdfExpandLabel(serverApplicationSecret, "iv", []byte{}, ivLength)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed deriving server application IV: %w", err)
	}

	return clientApplicationKey, clientApplicationIV, serverApplicationKey, serverApplicationIV, nil
}

func sendChangeCipherSuite(socketFD int) error {
	if err := syscall.Write(socketFD, []byte{0x14, 0x03, 0x03, 0x00, 0x01, 0x01}); err != nil {
		return fmt.Errorf("failed to write to socket: %w", err)
	}

	return nil
}

func readCipherSpecFromServer(socketFD int) error {
	changeCipherSpec, err := read(socketFD)
	if err != nil {
		return fmt.Errorf("failed reading from socket: %w", err)
	}

	if changeCipherSpec[0] != messageTypeCipherSpec {
		return NewError("changing cipher spec expected")
	}

	return nil
}

func clientVerificationData(clientHandshakeSecret, handshakeMessages []byte) ([]byte, error) {
	finishedKey, err := hkdfExpandLabel(clientHandshakeSecret, "finished", []byte{}, sha384Length)
	if err != nil {
		return nil, fmt.Errorf("failed expanding finished key label: %w", err)
	}

	hash := sha512.Sum384(handshakeMessages)
	hm := hmac.New(sha512.New384, finishedKey)

	hm.Write(hash[:])

	return hm.Sum(nil), nil
}

func clientHandshakeFinishedMessage(
	clientHandshakeSecret, handshakeMessages, clientHandshakeKey, clientHandshakeIV []byte,
) ([]byte, error) {
	verificationData, err := clientVerificationData(clientHandshakeSecret, handshakeMessages)
	if err != nil {
		return nil, fmt.Errorf("failed getting client verification data: %w", err)
	}

	verificationData = append(verificationData, messageTypeHandshake)

	message, err := binary.Append([]byte{0x14, 0x00, 0x00, 0x30}, binary.BigEndian, verificationData)
	if err != nil {
		return nil, fmt.Errorf("failed appending client verification data: %w", err)
	}

	messageLen, err := cast.ToUint16E(len(message) + authTagLength)
	if err != nil {
		return nil, fmt.Errorf("failed converting message length: %w", err)
	}

	additional := binary.BigEndian.AppendUint16([]byte{0x17, 0x03, 0x03}, messageLen)

	encryptedMsg, err := encrypt(clientHandshakeKey, clientHandshakeIV, message, additional)
	if err != nil {
		return nil, fmt.Errorf("failed encrypting client handshake message: %w", err)
	}

	return encryptedMsg, nil
}
