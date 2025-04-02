package tls

import (
	"fmt"

	"github.com/r-erema/go_sendbox/utils/os/syscall"
)

func HandshakeClientSide(socketFD int, hosts []string, pubKey [32]byte) error {
	secret, err := Rand32Bytes()
	if err != nil {
		return fmt.Errorf("failed generating secret: %w", err)
	}

	clientHelloMsg, err := encodeClientHello(hosts, []publicKey{{
		payload:       pubKey,
		exchangeGroup: x25519,
	}}, secret)
	if err != nil {
		return fmt.Errorf("failed encoding client hello message: %w", err)
	}

	if err = syscall.Write(socketFD, clientHelloMsg); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	serverHello, err := syscall.Read(socketFD)
	if err != nil {
		return fmt.Errorf("failed to read connection of socket fd `%d`: %w", socketFD, err)
	}

	helloMsg, err := decodeRawHello(serverHello)
	if err != nil {
		return fmt.Errorf("failed decoding server hello message: %w", err)
	}
	_ = helloMsg

	return nil
}

func HandshakeServerSide(socketFD int, pubKey [32]byte) error {
	secret, err := Rand32Bytes()
	if err != nil {
		return fmt.Errorf("failed generating secret: %w", err)
	}

	buf, err := syscall.Read(socketFD)
	if err != nil {
		return fmt.Errorf("failed to read connection of socket fd `%d`: %w", socketFD, err)
	}

	clientHelloMsg, err := decodeRawHello(buf)
	if err != nil {
		return fmt.Errorf("failed decoding client hello message: %w", err)
	}

	serverHelloRaw, err := encodeServerHello(publicKey{
		payload:       pubKey,
		exchangeGroup: x25519,
	}, secret, clientHelloMsg.cipherSuites[0])
	if err != nil {
		return fmt.Errorf("failed encoding client hello message: %w", err)
	}

	if err = syscall.Write(socketFD, serverHelloRaw); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
