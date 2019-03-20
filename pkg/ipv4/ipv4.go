package ipv4

import (
	"encoding/binary"
	"net"

	"github.com/mas9612/nwspeaker/pkg/checksum"
	"github.com/mas9612/nwspeaker/pkg/ethernet"
	"github.com/mas9612/nwspeaker/pkg/iface"
	"github.com/pkg/errors"
)

// Header represents IPv4 header.
type Header struct {
	Version        uint8
	IHL            uint8
	TypeOfService  uint8
	TotalLength    uint16 // Header len + Data len
	Identification uint16
	Flags          uint8
	FlagmentOffset uint16
	TimeToLive     uint8
	Protocol       uint8
	HeaderChecksum uint16
	SrcAddress     net.IP
	DstAddress     net.IP
}

// TODO: support IPv4 option

// Encode returns byte-encoded data of IPv4 header.
func (h *Header) Encode() []byte {
	buffer := make([]byte, HeaderLen)

	buffer[0] = (h.Version << 4) | h.IHL
	buffer[1] = h.TypeOfService
	binary.BigEndian.PutUint16(buffer[2:], h.TotalLength)
	binary.BigEndian.PutUint16(buffer[4:], h.Identification)

	fOffset := make([]byte, 2)
	binary.BigEndian.PutUint16(fOffset, h.FlagmentOffset)
	buffer[6] = (h.Flags << 5) | ((fOffset[0]) & (0xff >> 3))
	buffer[7] = fOffset[1]

	buffer[8] = h.TimeToLive
	buffer[9] = h.Protocol

	copy(buffer[12:], h.SrcAddress.To4())
	copy(buffer[16:], h.DstAddress.To4())

	checksum := checksum.SumOfOnesComplement16(buffer)
	copy(buffer[10:], checksum)

	return buffer
}

// Packet represents IPv4 packet.
type Packet struct {
	Header
	Data []byte
}

// Encode returns byte-encoded data of IPv4 packet.
func (p *Packet) Encode() []byte {
	header := p.Header.Encode()
	buffer := make([]byte, len(header)+len(p.Data))
	copy(buffer[0:], header)
	copy(buffer[len(header):], p.Data)
	return buffer
}

// Option is option which is used to send IP packet.
type Option func(*config)

// SetDstMac sets the destination MAC address.
func SetDstMac(dst net.HardwareAddr) Option {
	return func(c *config) {
		c.DstMac = dst
	}
}

type config struct {
	DstMac net.HardwareAddr
}

// Send sends given packet data to dst.
// packet must not include IPv4 header.
func Send(outIfname string, dst net.IP, payload []byte, proto uint8, opts ...Option) error {
	c := config{}
	for _, o := range opts {
		o(&c)
	}

	src, err := iface.IPv4AddressByName(outIfname)
	if err != nil {
		return errors.Wrap(err, "failed to get source IP address")
	}

	hdr := Header{
		Version:     Version4,
		IHL:         HeaderLen / 4,
		TotalLength: uint16(HeaderLen + len(payload)),
		// TODO: set identification properly
		Identification: 0,
		TimeToLive:     DefaultTTL,
		Protocol:       ProtoICMP,
		SrcAddress:     src,
		DstAddress:     dst,
	}
	pkt := &Packet{
		Header: hdr,
		Data:   payload,
	}

	// TODO: if DstMAC option is empty, resolve destination mac address with ARP
	return ethernet.Send(outIfname, c.DstMac, pkt, ethernet.TypeIPv4)
}
