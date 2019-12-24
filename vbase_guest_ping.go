package goproxmox

import (
	"errors"
	"fmt"
	"strings"
)

func (v *vbaseimpl) GuestPing() (bool, error) {
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/ping", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
	if err != nil {
		return false, err
	} else if resp.StatusCode == 500 && strings.Contains(resp.RawResponse.Status, "not running") {
		return false, nil
	} else if resp.StatusCode >= 400 {
		return false, errors.New(resp.RawResponse.Status)
	} else {
		return true, nil
	}
}
