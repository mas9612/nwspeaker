package command

import (
	"fmt"
	"net"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jessevdk/go-flags"
	"github.com/mas9612/nwspeaker/pkg/icmp"
	"github.com/mas9612/nwspeaker/pkg/ipv4"
)

// ICMPCommand is a command to craft ICMP packet.
type ICMPCommand struct{}

// Help returns long-formt help text of ICMPCommand.
func (c *ICMPCommand) Help() string {
	helpText := `
Usage: craftpkt icmp [options]

  Craft ICMP packet.

Options:
  -i, --interface   Output interface.
  --src-mac         Source MAC address.
  --dst-mac         Destination MAC address.
  --src-ip          Source IP address.
  --dst-ip          Destination IP address.
  -t, --type        ICMP type code.
  -l, --list-types  Print supported ICMP type codes and exit.
`
	return strings.TrimSpace(helpText)
}

func craftICMPEchoRequest() (*icmp.Message, error) {
	header := icmp.Message{
		Type:     icmp.TypeEcho,
		Code:     0,
		Checksum: 0,
	}
	payload := &icmp.Echo{}
	_ = header
	_ = payload
	return &icmp.Message{}, nil
}

// Run runs ICMPCommand and returns exit status.
func (c *ICMPCommand) Run(args []string) int {
	var opts struct {
		Interface string `short:"i" long:"interface"`
		SrcMac    string `long:"src-mac"`
		DstMac    string `long:"dst-mac"`
		SrcIP     string `long:"src-ip"`
		DstIP     string `long:"dst-ip"`
		Type      int    `short:"t" long:"type"`
		ListTypes bool   `short:"l" long:"list-types"`
	}
	if _, err := flags.ParseArgs(&opts, args); err != nil {
		return 1
	}

	if opts.ListTypes {
		printSupportedTypes()
		return 0
	}

	// TODO: check required flags

	echo, err := icmp.NewEcho(opts.Interface, opts.DstIP, opts.DstMac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create ICMP Echo packet: %v\n", err)
		return 1
	}

	dstIP := net.ParseIP(opts.DstIP)
	if dstIP == nil {
		fmt.Fprintf(os.Stderr, "failed to parse destination IP address\n")
		return 1
	}
	dstMac, err := net.ParseMAC(opts.DstMac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse destination MAC address\n")
		return 1
	}
	if err := ipv4.Send(opts.Interface, dstIP, echo.Encode(), ipv4.ProtoICMP, ipv4.SetDstMac(dstMac)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send ICMP Echo: %v\n", err)
		return 1
	}

	return 0
}

func printSupportedTypes() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer writer.Flush()

	fmt.Fprintf(writer, "TypeCode\tMessage\n")
	fmt.Fprintf(writer, "8\tEcho Request\n")
	fmt.Fprintf(writer, "0\tEcho Reply\n")
}

// Synopsis returns one-line synopsis of ICMPCommand.
func (c *ICMPCommand) Synopsis() string {
	return "Craft arbitrary ICMP packet."
}
