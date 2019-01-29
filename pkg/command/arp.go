package command

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/mas9612/nwspeaker/pkg/arp"
	"github.com/mas9612/nwspeaker/pkg/ethernet"
	"github.com/mas9612/nwspeaker/pkg/iface"
	"github.com/pkg/errors"
)

var (
	out    string
	srcMac string
	srcIP  string
	dstMac string
	dstIP  string
	op     string
)

// ArpCommand is a command to craft ARP packet.
type ArpCommand struct{}

// Help returns long-form help text of ArpCommand.
func (c *ArpCommand) Help() string {
	helpText := `
Usage: craftpkt arp [options]

  Craft ARP packet.
  Both ARP request and ARP reply can be crafted with this command.

Options:
  -interface  Network interface name which ARP packet will be sent from.
              Required.
  -src-mac    Source MAC address.
  -src-ip     Source IP address.
  -dst-mac    Destination MAC address. Ignored when -op is "request".
  -dst-ip     Destination IP address. Required when you craft ARP request.
  -op         ARP operation type. Only "request" or "reply" will be accepted.
              Default: "request"
`
	return strings.TrimSpace(helpText)
}

func craftARPRequest() (*arp.Packet, error) {
	packet := &arp.Packet{
		HType:    arp.HardwareTypeEthernet,
		PType:    arp.ProtocolTypeIPv4,
		HLen:     ethernet.EtherLen,
		PLen:     net.IPv4len,
		Op:       arp.OpRequest,
		DstHAddr: ethernet.Zero, // in ARP request, destination MAC address is fixed to zero
	}

	if srcMac == "" {
		mac, err := iface.MACAddressByName(out)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get MAC address")
		}
		packet.SrcHAddr = mac
	} else {
		mac, err := net.ParseMAC(srcMac)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid MAC address '%s'", mac)
		}
		packet.SrcHAddr = mac
	}

	if srcIP == "" {
		ip, err := iface.IPv4AddressByName(out)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get IPv4 address")
		}
		if ip == nil {
			return nil, errors.Wrapf(err, "no IPv4 address is assigned to \"%s\"", out)
		}
		packet.SrcPAddr = ip
	} else {
		ip := net.ParseIP(srcIP)
		if ip == nil {
			return nil, errors.Errorf("invalid IPv4 address '%s'", srcIP)
		}
		packet.SrcPAddr = ip
	}

	if dstIP == "" {
		return nil, errors.New("-dst-ip is required when you craft ARP request")
	}
	ip := net.ParseIP(dstIP)
	if ip == nil {
		return nil, errors.Errorf("invalid IPv4 address '%s'", dstIP)
	}
	packet.DstPAddr = ip

	return packet, nil
}

// Run runs ArpCommand and returns exit status.
func (c *ArpCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("arp", flag.ExitOnError)
	flagSet.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", c.Help()) }
	flagSet.StringVar(&out, "interface", "", "")
	flagSet.StringVar(&srcMac, "src-mac", "", "")
	flagSet.StringVar(&srcIP, "src-ip", "", "")
	flagSet.StringVar(&dstMac, "dst-mac", "", "")
	flagSet.StringVar(&dstIP, "dst-ip", "", "")
	flagSet.StringVar(&op, "op", "request", "")
	flagSet.Parse(args)

	if out == "" {
		fmt.Fprintln(os.Stderr, "-interface is required")
		return 1
	}
	if op != "request" && op != "reply" {
		fmt.Fprintln(os.Stderr, "invalid op type. valid type: \"request\", \"reply\"")
		return 1
	}

	var packet *arp.Packet
	switch op {
	case "request":
		var err error
		packet, err = craftARPRequest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
		dstMac = "ff:ff:ff:ff:ff:ff"
	case "reply":
	}

	if err := ethernet.Send(out, dstMac, packet, ethernet.TypeARP); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

// Synopsis returns one-line synopsis of ArpCommand.
func (c *ArpCommand) Synopsis() string {
	return "Craft arbitrary ARP packet."
}
