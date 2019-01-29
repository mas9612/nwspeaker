package iface

import (
	"net"

	"github.com/pkg/errors"
)

// IPv4AddressByName returns IPv4 address of given network interface.
// If no IPv4 address is assigned to it, nil will be returned.
func IPv4AddressByName(iface string) (net.IP, error) {
	out, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get iface info")
	}

	addrs, err := out.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get addresses from iface")
	}
	var ip net.IP
	for _, addr := range addrs {
		switch a := addr.(type) {
		case *net.IPNet:
			if !a.IP.IsLoopback() && a.IP.To4() != nil {
				ip = a.IP
			}
		}
	}
	return ip, nil
}

// MACAddressByName returns MAC address assigned to given network interface.
func MACAddressByName(iface string) (net.HardwareAddr, error) {
	out, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get iface info")
	}
	return out.HardwareAddr, nil
}
