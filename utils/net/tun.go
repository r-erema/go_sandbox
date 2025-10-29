package net

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	gvisortun "gvisor.dev/gvisor/pkg/tcpip/link/tun"
)

const (
	ipVerLen = 4
	ipV4     = 4
	ipV6     = 6
)

func SetupTun(name string, ipAddr, peer string) (*net.Interface, error) {
	if out, err := exec.Command("ip", "tuntap", "add", "dev", name, "mode", "tun").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding tun device `%s` error: %w %s", name, err, out)
	}

	if out, err := exec.Command("ip", "link", "set", "dev", name, "up").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("enabling link `%s` error: %w %s", name, err, out)
	}

	if out, err := exec.Command("ip", "addr", "add", ipAddr, "peer", peer, "dev", name).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding addr `%s` to `%s` error: %w %s", ipAddr, name, err, out)
	}

	veth, err := net.InterfaceByName(name)
	if err != nil {
		return nil, fmt.Errorf("getting tun `%s` error: %w", name, err)
	}

	return veth, nil
}

func OpenTun(tun *net.Interface) (*os.File, error) {
	fd, err := gvisortun.Open(tun.Name)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %w", tun.Name, err)
	}

	return os.NewFile(uintptr(fd), tun.Name), nil
}

func ReadFromTun(ctx context.Context, tun *net.Interface) (chan gopacket.Packet, error) {
	tunFile, err := OpenTun(tun)
	if err != nil {
		return nil, fmt.Errorf("opening tun %s error: %w", tun.Name, err)
	}

	var (
		buf    = make([]byte, maxIPPacketLength)
		tunCh  = make(chan gopacket.Packet)
		packet gopacket.Packet
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(tunCh)

				if err = tunFile.Close(); err != nil {
					log.Printf("close tun error: %v", err)
				}

				return
			default:
				n, err := tunFile.Read(buf)
				if err != nil {
					log.Printf("failed to read from tun file: %v", err)
				}

				buf = buf[:n]

				switch buf[0] >> ipVerLen {
				case ipV4:
					packet = gopacket.NewPacket(buf, layers.LayerTypeIPv4, gopacket.Default)
				case ipV6:
					packet = gopacket.NewPacket(buf, layers.LayerTypeIPv6, gopacket.Default)
				default:
					log.Printf("unknown packet type: %d", buf[0]>>ipV4)
				}

				tunCh <- packet
			}
		}
	}()

	return tunCh, nil
}
