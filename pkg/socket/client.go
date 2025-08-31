package socket

import (
	"fmt"
	"net"
	"slices"

	"github.com/r-erema/go_sendbox/pkg/syscall"
	"github.com/r-erema/go_sendbox/pkg/tls"
	"golang.org/x/sys/unix"
)

func ClientCall(domain string, port uint16, message []byte) ([]byte, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, fmt.Errorf("failed resolving domain: %w", err)
	}

	i := slices.IndexFunc(ips, func(item net.IP) bool { return item.To4() != nil })
	ip := ips[i].String()

	socketFD, err := syscall.SocketFD(unix.AF_INET, unix.SOCK_STREAM)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %w", err)
	}

	if err = syscall.Connect(socketFD, ip, port); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	privateKey, pubKey, err := tls.GeneratePrivateAndPublicKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private and public keys: %w", err)
	}

	clientAppKey, clientAppIV, serverAppKey, serverAppIV, err := tls.HandshakeClientSide(
		socketFD,
		[]string{domain},
		privateKey,
		pubKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to handshake client side: %w", err)
	}

	if err = tls.WriteAppData(socketFD, message, clientAppKey, tls.IV(0, clientAppIV)); err != nil {
		return nil, fmt.Errorf("failed to write client app data: %w", err)
	}

	// no idea so far what is that
	_, err = tls.ReadAppData(socketFD, serverAppKey, tls.IV(0, serverAppIV))
	if err != nil {
		return nil, fmt.Errorf("failed to read server app data: %w", err)
	}

	response, err := tls.ReadAppData(socketFD, serverAppKey, tls.IV(1, serverAppIV))
	if err != nil {
		return nil, fmt.Errorf("failed to read server app data: %w", err)
	}

	return response, nil
}
