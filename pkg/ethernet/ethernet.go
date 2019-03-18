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

// Option is option which is used to send ethernet frame.
type Option func(*config)

type config struct {
	srcMac net.HardwareAddr
}

// Send sends ethernet packet to given dst with given payload
func Send(outIfname string, dst net.HardwareAddr, payload Payload, proto uint16, opts ...Option) error {
	c := config{}
	for _, o := range opts {
		o(&c)
	}

	oif, err := net.InterfaceByName(outIfname)
	if err != nil {
		return errors.Wrap(err, "failed to get out iface info")
	}
	if c.srcMac == nil {
		c.srcMac = oif.HardwareAddr
	}

	header := Header{
		SrcAddr:   c.srcMac,
		DstAddr:   dst,
		EtherType: proto,
	}

	rawPayload := payload.Encode()
	packetLen := HeaderLen + len(rawPayload)
	packet := make([]byte, packetLen)
	copy(packet[0:], header.Encode())
	copy(packet[HeaderLen:], payload.Encode())

	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(proto))
	if err != nil {
		return errors.Wrap(err, "failed to open socket")
	}
	addr := &unix.SockaddrLinklayer{
		Protocol: proto,
		Ifindex:  oif.Index,
		Halen:    EtherLen,
	}
	copy(addr.Addr[:], dst[:])

	err = unix.Sendto(fd, packet, 0, addr)
	if err != nil {
		return errors.Wrap(err, "failed to send data")
	}

	return nil
}
