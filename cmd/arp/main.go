package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/mas9612/nwspeaker/pkg/arp"
	"github.com/mas9612/nwspeaker/pkg/ethernet"
)

var (
	iface  = flag.String("interface", "", "Interface name to send ARP packet.")
	target = flag.String("target", "", "Target IP address. For ARP request: IP address that you want to know MAC address. For ARP reply: IP address that will be sent ARP reply.")
)

func init() {
	flag.Parse()

	if *iface == "" {
		fmt.Fprintln(os.Stderr, "Please pass -interface flag")
		os.Exit(1)
	}
	if *target == "" {
		fmt.Fprintln(os.Stderr, "Please pass -target flag")
		os.Exit(1)
	}
}

func main() {
	out, err := net.InterfaceByName(*iface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get iface info: %v\n", err)
		os.Exit(1)
	}

	addrs, err := out.Addrs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get addresses from out iface: %v\n", err)
		os.Exit(1)
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

	payload, err := arp.NewRequest(*target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build ARP request packet\n")
		os.Exit(1)
	}
	payload.SrcHAddr = out.HardwareAddr
	payload.SrcPAddr = ip

	if err := ethernet.Send(*iface, "ff:ff:ff:ff:ff:ff", payload, ethernet.TypeARP); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
