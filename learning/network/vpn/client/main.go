package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	utilsnet "github.com/r-erema/go_sendbox/utils/net"
)

func main() {
	serverAddr := os.Args[1]
	tunName := os.Args[2]
	destIPsV4Input := os.Args[3]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	tun, err := utilsnet.SetupTun(tunName, "10.0.5.1", "10.0.5.2")
	if err != nil {
		log.Printf("[Client] Setup tun error: %v", err)
	}

	defer func() {
		if err = utilsnet.DeleteLink(tun.Name); err != nil {
			log.Printf("delete tun %s error: %s", tun.Name, err)
		}
	}()

	cidrsV4 := []net.IP{net.ParseIP(destIPsV4Input)}

	if err = utilsnet.DirectTrafficFromCidrsToDevice(tun.Name, cidrsV4, nil); err != nil {
		log.Printf("[Client] Direct traffic to device error: %v", err)

		return
	}

	defer func() {
		if err = utilsnet.StopCidrsTrafficFromDevice(tun.Name, cidrsV4, nil); err != nil {
			log.Printf("[Client] Cancel traffic from device error: %v", err)
		}
	}()

	udpConn, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Printf("[Client] Setup udp connection error: %v", err)

		return
	}

	defer func() {
		if err = udpConn.Close(); err != nil {
			log.Printf("[Client] Close udp connection error: %s", err)
		}
	}()

	log.Printf("[Client] Successfully connected to UDP Server %s", serverAddr)

	packetsCh, err := utilsnet.ReadFromTun(ctx, tun)
	if err != nil {
		log.Printf("[Client] Failed to read packet from server: %v", err)

		return
	}

	for packet := range packetsCh {
		log.Printf(
			"[Client] Source IP: %s Destination IP: %s\n",
			packet.NetworkLayer().NetworkFlow().Src().String(),
			packet.NetworkLayer().NetworkFlow().Dst().String(),
		)

		if _, err = fmt.Fprintln(udpConn, string(packet.Data())); err != nil {
			log.Printf("[Client] Failed to send packet to server: %v", err)
		}

		log.Printf("[Client] Sent packet to Server %s\n", udpConn.RemoteAddr().String())
	}
}
