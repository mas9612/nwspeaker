package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"

	"github.com/jessevdk/go-flags"
	"github.com/mas9612/nwspeaker/pkg/arp"
	"github.com/mas9612/nwspeaker/pkg/endian"
	"github.com/mas9612/nwspeaker/pkg/ethernet"
	"github.com/mas9612/nwspeaker/pkg/iface"
)

type options struct {
	Interface string `short:"i" long:"interface" required:"true" description:"Output interface name. Required."`
	Garp      bool   `short:"g" long:"garp" description:"Send GARP instead of normal ARP request."`
	Check     bool   `short:"c" long:"check" description:"Check the response from other host. If this is not true, simply send data and exit."`
	Args      struct {
		TargetIP string `description:"IP address want to get MAC address. Not used when -g flag is on."`
	} `positional-args:"yes"`
}

func main() {
	opts := options{}
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	if !opts.Garp && opts.Args.TargetIP == "" {
		fmt.Fprintf(os.Stderr, "TargetIP is required when -g flag is not set.\n")
		os.Exit(1)
	}

	oif, err := net.InterfaceByName(opts.Interface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get interface information: %s\n", err.Error())
		os.Exit(1)
	}
	myip, err := iface.IPv4AddressByName(opts.Interface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get ipv4 address: %s\n", err.Error())
		os.Exit(1)
	}

	var data *arp.Packet
	if opts.Garp {
		var err error
		data, err = arp.NewRequest(myip.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create ARP packet: %s\n", err.Error())
			os.Exit(1)
		}
		// according to RFC5227, we should only set sender hardware address and target ip address when send GARP.
		// https://tools.ietf.org/html/rfc5227#section-2.1.1
		data.SrcHAddr = oif.HardwareAddr
	} else {
		var err error
		data, err = arp.NewRequest(opts.Args.TargetIP)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create ARP packet: %s\n", err.Error())
			os.Exit(1)
		}
		data.SrcHAddr = oif.HardwareAddr
		data.SrcPAddr = myip
	}

	// prepare raw socket
	soc, err := ethernet.Dial(endian.Htons(ethernet.TypeARP))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create ethernet raw socket: %s\n", err.Error())
		os.Exit(1)
	}
	addr := &unix.SockaddrLinklayer{
		Protocol: endian.Htons(ethernet.TypeARP),
		Ifindex:  oif.Index,
		Halen:    ethernet.EtherLen,
	}
	if err := soc.Bind(addr); err != nil {
		fmt.Fprintf(os.Stderr, "failed to bind interface: %s\n", err.Error())
		os.Exit(1)
	}

	if err := soc.Send(data, 0, "ff:ff:ff:ff:ff:ff"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send ARP frame: %s\n", err.Error())
		os.Exit(1)
	}
	if opts.Check {
		b, err := soc.Recv(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to receive ARP reply: %s\n", err.Error())
			os.Exit(1)
		}
		res := arp.Parse(b)
		if res.SrcPAddr.String() == opts.Args.TargetIP {
			fmt.Printf("MAC address of %s is %s\n", res.SrcPAddr.String(), res.SrcHAddr.String())
		} else {
			fmt.Printf("could not get the MAC address\n")
		}
	}
}
