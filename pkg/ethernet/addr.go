package ethernet

import "net"

var (
	// Broadcast represents the ethernet broadcast address
	Broadcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)
