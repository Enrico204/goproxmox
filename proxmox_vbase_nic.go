package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"strings"
)

type VBaseNICSettings struct {
	Id     int
	Bridge string
	Tag    *int

	// Virtual machines
	Model *string

	// LXC
	Name     string
	DHCP     bool
	IPv4     *net.IPNet
	Gateway4 *net.IPAddr
}

func (v *vbaseimpl) SetNIC(settings *VBaseNICSettings) error {
	niccfg := []string{
		"bridge=" + settings.Bridge,
	}
	if settings.Tag != nil {
		niccfg = append(niccfg, fmt.Sprintf("tag=%d", *settings.Tag))
	}

	if v.vmtype == "lxc" {
		niccfg = append(niccfg, "name="+settings.Name)
		if settings.DHCP {
			niccfg = append(niccfg, "ip=dhcp")
		} else if settings.IPv4 != nil {
			niccfg = append(niccfg, "ip="+settings.IPv4.String())
		}

		if settings.Gateway4 != nil {
			niccfg = append(niccfg, "gw="+settings.Gateway4.String())
		}
	} else {
		model := "virtio"
		if settings.Model != nil {
			model = *settings.Model
		}
		niccfg = append(niccfg, "model="+model)
	}

	bodyparams := map[string]string{
		fmt.Sprintf("net%d", settings.Id): strings.Join(niccfg, ","),
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
