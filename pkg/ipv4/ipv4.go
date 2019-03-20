package ipv4

import (
	"encoding/binary"
	"net"

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

	checksum := calculateChecksum(buffer)
	copy(buffer[10:], checksum)

	return buffer
}

func onesComplement(b byte) byte {
	var complement byte
	for i := 7; i >= 0; i-- {
		if (b >> uint(i) & 0x1) == 0x0 {
			complement |= 0x1
		}
		if i > 0 {
			complement <<= 1
		}
	}
	return complement
}

func bytesOnesComplement(data []byte) []byte {
	complement := make([]byte, len(data))
	for i, d := range data {
		complement[i] = onesComplement(d)
	}
	return complement
}

func calculateChecksum(b []byte) []byte {
	var sum uint32
	for i := 0; i < len(b); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(b[i : i+2]))
		if (sum >> 16) > 0x0 {
			sum += sum >> 16
			sum &= 0xffff // clear carry bit
		}
	}
	sumBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sumBytes, uint16(sum))
	return bytesOnesComplement(sumBytes)
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

type config struct{}

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

	dstMac, err := iface.MACAddressByName(outIfname)
	if err != nil {
		return errors.Wrap(err, "failed to get destination MAC address")
	}
	return ethernet.Send(outIfname, dstMac, pkt, ethernet.TypeIPv4)
}
