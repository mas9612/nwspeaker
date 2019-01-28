package ethernet

import "net"

var (
	// Zero represents the zero address
	Zero = net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// Broadcast represents the ethernet broadcast address
	Broadcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)
