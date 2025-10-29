//nolint:gosec // disable G204
package net

import (
	"errors"
	"fmt"
	"net"
	"os/exec"

	"k8s.io/apimachinery/pkg/util/json"
)

const maxIPPacketLength = 65535

type IPAddr struct {
	Ifindex   int      `json:"ifindex"`
	Ifname    string   `json:"ifname"`
	Flags     []string `json:"flags"`
	MTU       int      `json:"mtu"`
	Qdisc     string   `json:"qdisc"`
	Operstate string   `json:"operstate"`
	Group     string   `json:"group"`
	Txqlen    int      `json:"txqlen"`
	LinkType  string   `json:"link_type"`
	Address   string   `json:"address"`
	Broadcast string   `json:"broadcast"`
	AddrInfo  []struct {
		Family            string `json:"family"`
		Local             string `json:"local"`
		Prefixlen         int    `json:"prefixlen"`
		Scope             string `json:"scope"`
		Label             string `json:"label,omitempty"`
		ValidLifeTime     int64  `json:"valid_life_time"`
		PreferredLifeTime int64  `json:"preferred_life_time"`
	} `json:"addr_info"`
}

var (
	ErrIPAddrNotFound     = errors.New("ip address not found")
	ErrIPAddrAlreadyInUse = errors.New("ip address already in use")
)

func IPAddrList() ([]*IPAddr, error) {
	cmd := exec.Command("ip", "--json", "addr", "show")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("showing IP address list error: %w", err)
	}

	var IPs []*IPAddr

	err = json.Unmarshal(out, &IPs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON output: %w", err)
	}

	return IPs, nil
}

func IsIPAddrAlreadyInUse(ip string) (bool, error) {
	ipAddr, cidr, err := net.ParseCIDR(ip)
	if err != nil {
		return false, fmt.Errorf("parsing CIDR error: %w", err)
	}

	addrs, err := IPAddrList()
	if err != nil {
		return false, fmt.Errorf("getting addresses list error: %w", err)
	}

	var leadingOnesMaskSize int

	for _, addr := range addrs {
		for _, info := range addr.AddrInfo {
			leadingOnesMaskSize, _ = cidr.Mask.Size()
			if info.Local == ipAddr.String() && info.Prefixlen == leadingOnesMaskSize {
				return true, nil
			}
		}
	}

	return false, nil
}

func SetupLoopBackInterface(namespace *string) (*net.Interface, error) {
	cmd := exec.Command("ip", "link", "set", "dev", "lo", "up")

	if namespace != nil {
		cmd = exec.Command("ip", "netns", "exec", *namespace, "ip", "link", "set", "dev", "lo", "up")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("adding loopback interface error: %w %s", err, out)
	}

	lo, err := net.InterfaceByName("lo")
	if err != nil {
		return nil, fmt.Errorf("getting loopback interface error: %w", err)
	}

	return lo, nil
}

func CreateNS(nsName string) error {
	if out, err := exec.Command("ip", "netns", "add", nsName).CombinedOutput(); err != nil {
		return fmt.Errorf("creating netns %s error: %w %s", nsName, err, out)
	}

	if _, err := SetupLoopBackInterface(&nsName); err != nil {
		return fmt.Errorf("setting up loopback interface error: %w", err)
	}

	return nil
}

func DeleteNS(nsName string) error {
	if out, err := exec.Command("ip", "netns", "delete", nsName).CombinedOutput(); err != nil {
		return fmt.Errorf("deleting netns %s error: %w %s", nsName, err, out)
	}

	return nil
}

func SetupVeth(vethName, vethCIDR, peerVethName, peerAddr, peerNSName string) (*net.Interface, error) {
	if out, err := exec.Command(
		"ip", "link", "add", vethName, "type", "veth",
		"peer", "name", peerVethName, "netns", peerNSName).
		CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding veth `%s - %s` error: %w %s", vethName, peerVethName, err, out)
	}

	if out, err := exec.Command("ip", "addr", "add", vethCIDR, "dev", vethName).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding IP address to veth `%s` error: %w %s", vethName, err, out)
	}

	if out, err := exec.Command("ip", "link", "set", "dev", vethName, "up").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("enabling link `%s` error: %w %s", vethName, err, out)
	}

	if out, err := exec.Command("ip", "netns", "exec", peerNSName, "ip", "addr", "add", peerAddr, "dev",
		peerVethName).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("adding IP address to veth `%s` error: %w %s", peerVethName, err, out)
	}

	if out, err := exec.Command("ip", "netns", "exec", peerNSName, "ip", "link", "set", "dev",
		peerVethName, "up").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("enabling link `%s` error: %w %s", peerVethName, err, out)
	}

	if out, err := exec.Command("ip", "netns", "exec", peerNSName, "ip", "link", "set", "dev",
		peerVethName, "up").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("enabling link `%s` error: %w %s", peerVethName, err, out)
	}

	veth, err := net.InterfaceByName(vethName)
	if err != nil {
		return nil, fmt.Errorf("getting veth `%s` error: %w", vethName, err)
	}

	return veth, nil
}

func SetupBridge(name, ipAddr string) error {
	foundBridge, err := net.InterfaceByName(name)
	if err != nil && err.Error() != "route ip+net: no such network interface" {
		return fmt.Errorf("checking existed bridge `%s` error: %w", name, err)
	}

	if foundBridge == nil {
		cmd := exec.Command("ip", "link", "add", "name", name, "type", "bridge")

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("adding bridge `%s` error: %w", name, err)
		}
	}

	if err = EnableDevice(name); err != nil {
		return fmt.Errorf("enabling bridge `%s` error: %w", name, err)
	}

	addrAlreadyInUse, err := IsIPAddrAlreadyInUse(ipAddr)
	if err != nil && !errors.Is(err, ErrIPAddrNotFound) {
		return fmt.Errorf("checking existed address `%s` error: %w", name, err)
	}

	if !addrAlreadyInUse {
		cmd := exec.Command("ip", "addr", "add", ipAddr, "dev", name)

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("attaching IP `%s` to the bridge `%s`: %w", ipAddr, name, err)
		}
	} else {
		return ErrIPAddrAlreadyInUse
	}

	return nil
}

func RemoveBridge(name string) error {
	cmd := exec.Command("ip", "link", "delete", name, "type", "bridge")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("removal bridge `%s` error: %w %s", name, err, out)
	}

	return nil
}

func AttachDeviceToBridge(deviceName, bridgeName string) error {
	cmd := exec.Command("ip", "link", "set", "dev", deviceName, "master", bridgeName)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"attaching device `%s` to the bridge `%s` error: %w %s",
			deviceName,
			bridgeName,
			err,
			out,
		)
	}

	return nil
}

func AddIPAddrToInterface(ip, interfaceName string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", interfaceName)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"adding IP `%s` to the interface `%s` error: %w %s",
			ip,
			interfaceName,
			err,
			out,
		)
	}

	return nil
}

func SetDefaultGateway(ip string) error {
	if out, err := exec.Command("ip", "route", "add", "default", "via", ip).CombinedOutput(); err != nil {
		return fmt.Errorf("adding default route via ip `%s` error: %w %s", ip, err, out)
	}

	return nil
}

func EnableDevice(deviceName string) error {
	if out, err := exec.Command("ip", "link", "set", "dev", deviceName, "up").CombinedOutput(); err != nil {
		return fmt.Errorf("enabling device `%s` error: %w %s", deviceName, err, out)
	}

	return nil
}

func DeleteLink(linkName string) error {
	cmd := exec.Command("ip", "link", "delete", linkName)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runnning `%s` command error: %w %s", cmd.String(), err, out)
	}

	return nil
}
