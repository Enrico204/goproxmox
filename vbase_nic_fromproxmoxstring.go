package goproxmox

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func (settings *VBaseNICSettings) FromProxmoxString(cfg string) error {
	tags := strings.Split(cfg, ",")
	for _, t := range tags {
		s := strings.SplitN(t, "=", 2)
		switch s[0] {
		case "bridge":
			settings.Bridge = s[1]
		case "tag":
			vlanTag, err := strconv.ParseUint(s[1], 10, 32)
			if err != nil {
				return err
			}
			settings.Tag = uint(vlanTag)
		case "firewall":
			settings.Firewall = s[1] == "1"
		case "link_down":
			settings.LinkDown = s[1] == "1"
		case "rate":
			rateMbps, err := strconv.ParseFloat(s[1], 32)
			if err != nil {
				return err
			}
			settings.Rate = float32(rateMbps)
		case "trunks":
			settings.Trunks = strings.Split(s[1], ",")
		case "name":
			settings.Name = s[1]
		case "mtu":
			mtuSize, err := strconv.ParseUint(s[1], 10, 32)
			if err != nil {
				return err
			}
			settings.MTU = uint(mtuSize)
		case "ip":
			if s[1] == "manual" {
				settings.Manualv4 = true
			} else if s[1] == "dhcp" {
				settings.DHCPv4 = true
			} else {
				ip, ipnet, err := net.ParseCIDR(s[1])
				if err != nil {
					return err
				}
				settings.IPv4 = ip
				settings.Maskv4 = ipnet.Mask
			}
		case "ip6":
			if s[1] == "manual" {
				settings.Manualv6 = true
			} else if s[1] == "dhcp" {
				settings.DHCPv6 = true
			} else if s[1] == "auto" {
				settings.Autov6 = true
			} else {
				ip, ipnet, err := net.ParseCIDR(s[1])
				if err != nil {
					return err
				}
				settings.IPv6 = ip
				settings.Maskv6 = ipnet.Mask
			}
		case "gw":
			settings.Gateway4 = net.ParseIP(s[1])
		case "gw6":
			settings.Gateway6 = net.ParseIP(s[1])
		case "macaddr":
			fallthrough
		case "e1000-82540em":
			fallthrough
		case "e1000-82544gc":
			fallthrough
		case "e1000-82545em":
			fallthrough
		case "i82551":
			fallthrough
		case "model":
			fallthrough
		case "i82557b":
			fallthrough
		case "i82559er":
			fallthrough
		case "rtl8139":
			fallthrough
		case "ne2k_pci":
			fallthrough
		case "e1000":
			fallthrough
		case "pcnet":
			fallthrough
		case "virtio":
			fallthrough
		case "vmxnet3":
			fallthrough
		case "ne2k_isa":
			hwaddr, err := net.ParseMAC(s[1])
			if err != nil {
				return err
			}
			settings.HardwareAddress = &hwaddr
			settings.Model = s[0]
		default:
			return errors.New("unhandled parameter in network config: " + t)
		}
	}
	return nil
}
