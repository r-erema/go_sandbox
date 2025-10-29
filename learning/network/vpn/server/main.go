package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	bufSize  = 128
	ipVerLen = 4
)

func main() {
	serverAddr := os.Args[1]

	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("[Server] Failed to resolve udp address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("[Server] Failed to listen UDP addr %v: %v", addr, err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatalf("[Server] Close udp error: %s", err)
		}
	}()

	buf := make([]byte, bufSize)

	var packet gopacket.Packet

	ctx := context.Background()

	log.Printf("[Server] Started at %s\n", conn.LocalAddr().String())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("[Server] Failed to read from UDP: %v", err)

				continue
			}

			buf = buf[:n]

			switch buf[0] >> ipVerLen {
			case net.IPv4len:
				packet = gopacket.NewPacket(buf, layers.LayerTypeIPv4, gopacket.Default)

				log.Printf(
					"[Server] Source IP: %s Destination IP: %s\n",
					packet.NetworkLayer().NetworkFlow().Src().String(),
					packet.NetworkLayer().NetworkFlow().Dst().String(),
				)

				// forward packet to the destination
			default:
				log.Printf("[Server] Unknown packet type: %d", buf[0]>>net.IPv4len)
			}
		}
	}
}
