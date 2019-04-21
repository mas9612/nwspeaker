package arp

import (
	"encoding/binary"
	"net"

	"github.com/mas9612/nwspeaker/pkg/ethernet"
	"github.com/pkg/errors"
)

// Packet represents ARP packet format.
type Packet struct {
	HType    uint16
	PType    uint16
	HLen     uint8
	PLen     uint8
	Op       uint16
	SrcHAddr net.HardwareAddr
	SrcPAddr net.IP
	DstHAddr net.HardwareAddr
	DstPAddr net.IP
}

// NewRequest returns Packet struct initialized as ARP request packet.
func NewRequest(dst string) (*Packet, error) {
	dstIP := net.ParseIP(dst)
	if dstIP == nil {
		return nil, errors.Errorf("failed to parse given address '%s'", dst)
	}
	if dstIP.To4() == nil {
		return nil, errors.Errorf("given address '%s' is not an IPv4 address", dst)
	}

	return &Packet{
		HType:    HardwareTypeEthernet,
		PType:    ProtocolTypeIPv4,
		HLen:     ethernet.EtherLen,
		PLen:     net.IPv4len,
		Op:       OpRequest,
		DstHAddr: ethernet.Zero,
		DstPAddr: dstIP,
	}, nil
}

// NewReply returns Packet struct initialized as ARP reply packet.
func NewReply(dstHAddr, dstPAddr string) (*Packet, error) {
	dstMAC, err := net.ParseMAC(dstHAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse MAC address '%s'", dstHAddr)
	}
	dstIP := net.ParseIP(dstPAddr)
	if dstIP == nil {
		return nil, errors.Errorf("failed to parse IP address '%s'", dstPAddr)
	}

	return &Packet{
		HType:    HardwareTypeEthernet,
		PType:    ProtocolTypeIPv4,
		HLen:     ethernet.EtherLen,
		PLen:     net.IPv4len,
		Op:       OpReply,
		DstHAddr: dstMAC,
		DstPAddr: dstIP,
	}, nil
}

// Encode returns byte-encoded data to send ARP packet to network.
func (p *Packet) Encode() []byte {
	payload := make([]byte, payloadLen)

	binary.BigEndian.PutUint16(payload[0:], p.HType)
	binary.BigEndian.PutUint16(payload[2:], p.PType)
	payload[4] = p.HLen
	payload[5] = p.PLen
	binary.BigEndian.PutUint16(payload[6:], p.Op)
	copy(payload[8:], p.SrcHAddr)
	copy(payload[14:], p.SrcPAddr.To4())
	copy(payload[18:], p.DstHAddr)
	copy(payload[24:], p.DstPAddr.To4())

	return payload
}

// Parse parses given frame and returns a pointer to Packet instance.
func Parse(b []byte) *Packet {
	_ = ethernet.Parse(b)

	b = b[ethernet.HeaderLen:]
	if len(b) < payloadLen { // incomplete frame
		return nil
	}

	p := &Packet{
		SrcHAddr: make([]byte, ethernet.EtherLen),
		DstHAddr: make([]byte, ethernet.EtherLen),
	}
	p.HType = binary.BigEndian.Uint16(b)
	p.PType = binary.BigEndian.Uint16(b[2:])
	p.HLen = b[4]
	p.PLen = b[5]
	p.Op = binary.BigEndian.Uint16(b[6:])
	copy(p.SrcHAddr, b[8:])
	p.SrcPAddr = net.IPv4(b[14], b[15], b[16], b[17])
	copy(p.DstHAddr, b[18:])
	p.DstPAddr = net.IPv4(b[24], b[25], b[26], b[27])
	return p
}
