package ethernet

import (
	"encoding/binary"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Header represents the ethernet header format.
type Header struct {
	SrcAddr   net.HardwareAddr
	DstAddr   net.HardwareAddr
	EtherType uint16
}

// Encode returns byte-encoded Header struct
func (h *Header) Encode() []byte {
	header := make([]byte, HeaderLen)
	copy(header[0:], h.SrcAddr)
	copy(header[6:], h.DstAddr)
	binary.BigEndian.PutUint16(header[12:], h.EtherType)
	return header
}

// Payload represents the application data of ethernet packet.
type Payload interface {
	Encode() []byte
}

// Send sends ethernet packet to given dst with given payload
func Send(outIfname, dst string, payload Payload, proto uint16) error {
	oif, err := net.InterfaceByName(outIfname)
	if err != nil {
		return errors.Wrap(err, "failed to get out iface info")
	}
	hw, err := net.ParseMAC(dst)
	if err != nil {
		return errors.Wrap(err, "faield to parse dst MAC address")
	}

	header := Header{
		SrcAddr:   oif.HardwareAddr,
		DstAddr:   hw,
		EtherType: proto,
	}

	rawPayload := payload.Encode()
	packetLen := HeaderLen + len(rawPayload)
	packet := make([]byte, packetLen)
	copy(packet[0:], header.Encode())
	copy(packet[HeaderLen:], payload.Encode())

	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, unix.ETH_P_ARP)
	if err != nil {
		return errors.Wrap(err, "failed to open socket")
	}
	addr := &unix.SockaddrLinklayer{
		Protocol: unix.ETH_P_ARP,
		Ifindex:  oif.Index,
		Halen:    EtherLen,
	}
	copy(addr.Addr[:], hw[:])

	err = unix.Sendto(fd, packet, 0, addr)
	if err != nil {
		return errors.Wrap(err, "failed to send data")
	}

	return nil
}
