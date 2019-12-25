package goproxmox

import (
	"net"
)

type VBaseNICSettings struct {
	// NIC Proxmox internal number
	Id int

	// Bridge where the NIC is attached
	Bridge string

	// VLAN ID. If zero, the tag is omitted (eg. untagged)
	Tag uint

	// Whether apply global firewall settings or not
	Firewall bool

	// Hardware address
	HardwareAddress *net.HardwareAddr

	// Rate limit the speed to this mbps
	Rate float32

	// VLAN trunks
	Trunks []string

	// *** Virtual machines only

	// NIC model (see Proxmox docs)
	Model string

	// Whether the link is down (eg. cable unplugged)
	LinkDown bool

	// *** LXC only
	Name string
	MTU  uint

	Manualv4 bool
	DHCPv4   bool
	IPv4     net.IP
	Maskv4   net.IPMask
	Gateway4 net.IP

	Manualv6 bool
	DHCPv6   bool
	Autov6   bool
	IPv6     net.IP
	Maskv6   net.IPMask
	Gateway6 net.IP
}
