package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"strings"
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
	Rate int

	// VLAN trunks
	Trunks []string

	// *** Virtual machines only

	// NIC model (see Proxmox docs)
	Model string

	// Whether the link is down (eg. cable unplugged)
	LinkDown bool

	// *** LXC only
	Name string
	MTU  int

	Manualv4 bool
	DHCPv4   bool
	IPv4     net.IPNet
	Gateway4 net.IP

	Manualv6 bool
	DHCPv6   bool
	Autov6   bool
	IPv6     net.IPNet
	Gateway6 net.IP
}

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

func (v *vbaseimpl) SetNIC(settings VBaseNICSettings) error {
	bodyparams := map[string]string{
		fmt.Sprintf("net%d", settings.Id): settings.ToProxmoxString(v.vmtype),
	}

	resp, err := v.node.proxmox.session.Put(
		fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/config", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		&grequests.RequestOptions{
			Data: bodyparams,
		})

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}

func (v *vbaseimpl) DeleteNIC(id int) error {
	resp, err := v.node.proxmox.session.Put(
		fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/config", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		&grequests.RequestOptions{
			Data: map[string]string{
				"delete": fmt.Sprintf("net%d", id),
			},
		})

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
