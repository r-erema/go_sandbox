//nolint:gosec // disable G204
package net

import (
	"fmt"
	"net"
	"os/exec"
)

func DirectTrafficFromCidrsToDevice(deviceName string, ipsV4, ipsV6 []net.IP) error {
	for i := range ipsV4 {
		if out, err := exec.Command("ip", "route", "add", ipsV4[i].String(), "dev", deviceName).
			CombinedOutput(); err != nil {
			return fmt.Errorf("adding route(%s) for device `%s` error: %w %s", ipsV4[i].String(), deviceName, err, out)
		}
	}

	for i := range ipsV6 {
		if out, err := exec.Command("ip", "-6", "route", "add", ipsV6[i].String(), "dev", deviceName).
			CombinedOutput(); err != nil {
			return fmt.Errorf("adding route(%s) for device `%s` error: %w %s", ipsV6[i].String(), deviceName, err, out)
		}
	}

	return nil
}

func StopCidrsTrafficFromDevice(deviceName string, cidrsV4, cidrsV6 []net.IP) error {
	for i := range cidrsV4 {
		if out, err := exec.Command("ip", "route", "del", cidrsV4[i].String(), "dev", deviceName).
			CombinedOutput(); err != nil {
			return fmt.Errorf("deleting route from device `%s` error: %w %s", deviceName, err, out)
		}
	}

	for i := range cidrsV6 {
		if out, err := exec.Command("ip", "-6", "route", "del", cidrsV6[i].String(), "dev", deviceName).
			CombinedOutput(); err != nil {
			return fmt.Errorf("deleting route from device `%s` error: %w %s", deviceName, err, out)
		}
	}

	return nil
}
