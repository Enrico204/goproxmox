package goproxmox

import (
	"fmt"
	"strings"
)

func (settings VBaseNICSettings) ToProxmoxString(vmtype string) string {
	niccfg := []string{
		"bridge=" + settings.Bridge,
	}
	if settings.Tag != 0 {
		niccfg = append(niccfg, fmt.Sprintf("tag=%d", settings.Tag))
	}
	if settings.Firewall {
		niccfg = append(niccfg, "firewall=1")
	} else {
		niccfg = append(niccfg, "firewall=0")
	}
	if settings.Rate != 0 {
		niccfg = append(niccfg, fmt.Sprintf("rate=%d", settings.Tag))
	}
	if len(settings.Trunks) > 0 {
		niccfg = append(niccfg, "trunks="+strings.Join(settings.Trunks, ";"))
	}

	if vmtype == "lxc" {
		if settings.HardwareAddress != nil {
			niccfg = append(niccfg, "hwaddr="+settings.HardwareAddress.String())
		}
		if settings.MTU > 0 {
			niccfg = append(niccfg, fmt.Sprintf("mtu=%d", settings.MTU))
		}

		niccfg = append(niccfg, "name="+settings.Name)
		if settings.DHCPv4 {
			niccfg = append(niccfg, "ip=dhcp")
		} else if settings.Manualv4 {
			niccfg = append(niccfg, "ip=manual")
		} else if settings.IPv4.String() != "" && !settings.IPv4.IP.IsUnspecified() {
			niccfg = append(niccfg, "ip="+settings.IPv4.String())
		}

		if !settings.Gateway4.Equal(nil) && !settings.Gateway4.IsUnspecified() {
			niccfg = append(niccfg, "gw="+settings.Gateway4.String())
		}

		// TODO: IPv6
	} else {
		model := "virtio"
		if settings.Model != "" {
			model = settings.Model
		}
		niccfg = append(niccfg, "model="+model)

		if settings.HardwareAddress != nil {
			niccfg = append(niccfg, "macaddr="+settings.HardwareAddress.String())
		}
	}

	return strings.Join(niccfg, ",")
}
