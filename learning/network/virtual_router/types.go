package virtual_router

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

const (
	ipLengthBytes  = 4
	macLengthBytes = 6

	messageLengthBytes = 1
)

type ExternalEthernetFrame struct {
	SourceMAC      net.HardwareAddr
	DestinationMAC net.HardwareAddr
	IpPacket       IpPacket
}

type IpPacket struct {
	SourceIP      net.IP
	DestinationIP net.IP
	Payload       []byte
}

func buildFrame(
	sourceMac, destinationMac net.HardwareAddr,
	sourceIP, destinationIP net.IP,
	payload []byte,
) ([]byte, error) {
	frame := &ExternalEthernetFrame{
		SourceMAC:      sourceMac,
		DestinationMAC: destinationMac,
		IpPacket: IpPacket{
			SourceIP:      sourceIP,
			DestinationIP: destinationIP,
			Payload:       payload,
		},
	}

	buf := new(bytes.Buffer)

	err := gob.NewEncoder(buf).Encode(frame)
	if err != nil {
		return nil, fmt.Errorf("failed to encode frame: %w", err)
	}

	return buf.Bytes(), nil
}
