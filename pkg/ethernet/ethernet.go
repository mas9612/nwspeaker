package ethernet

import (
	"encoding/binary"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Header represents the ethernet header format.
type Header struct {
	DstAddr   net.HardwareAddr
	SrcAddr   net.HardwareAddr
	EtherType uint16
}

// Encode returns byte-encoded Header struct
func (h *Header) Encode() []byte {
	header := make([]byte, HeaderLen)
	copy(header[0:], h.DstAddr)
	copy(header[6:], h.SrcAddr)
	binary.BigEndian.PutUint16(header[12:], h.EtherType)
	return header
}

// Payload represents the application data of ethernet packet.
type Payload interface {
	Encode() []byte
}

// Socket represents an ethernet socket used to send or receive data.
type Socket struct {
	fd    int
	proto uint16
	iface *net.Interface
}

// Dial returns new Socket instance.
func Dial(proto uint16) (*Socket, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(proto))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open raw socket")
	}
	return &Socket{
		fd:    fd,
		proto: proto,
	}, nil
}

// Bind binds interface to Socket instance.
func (s *Socket) Bind(sa unix.Sockaddr) error {
	if err := unix.Bind(s.fd, sa); err != nil {
		return errors.Wrap(err, "failed to bind interface to socket")
	}
	iface, err := net.InterfaceByIndex(sa.(*unix.SockaddrLinklayer).Ifindex)
	if err != nil {
		return errors.Wrap(err, "failed to get interface information")
	}
	s.iface = iface
	return nil
}

// Send sends given payload to dst.
func (s *Socket) Send(payload Payload, flags int, dst string) error {
	hw, err := net.ParseMAC(dst)
	if err != nil {
		return errors.Wrap(err, "failed to parse destination MAC address")
	}
	sa := &unix.SockaddrLinklayer{
		Protocol: s.proto,
		Ifindex:  s.iface.Index,
		Halen:    EtherLen,
	}
	copy(sa.Addr[:], hw)

	// add Ethernet header
	pb := payload.Encode()
	frame := make([]byte, HeaderLen+len(pb))
	copy(frame[HeaderLen:], pb)
	hdr := Header{
		SrcAddr:   s.iface.HardwareAddr,
		DstAddr:   hw,
		EtherType: s.proto,
	}
	copy(frame, hdr.Encode())

	if err := unix.Sendto(s.fd, frame, flags, sa); err != nil {
		return errors.Wrap(err, "send failed")
	}
	return nil
}

// Recv receives data from socket.
func (s *Socket) Recv(flags int) ([]byte, error) {
	buffer := make([]byte, BufferLen)
	_, _, err := unix.Recvfrom(s.fd, buffer, flags)
	if err != nil {
		return nil, errors.Wrap(err, "recv failed")
	}
	return buffer, nil
}

// Close closes socket.
func (s *Socket) Close() error {
	return unix.Close(s.fd)
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
