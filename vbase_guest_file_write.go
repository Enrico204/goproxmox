package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"gitlab.com/Enrico204/goproxmox/proxmoxssh"
	"strconv"
)

func (v *vbaseimpl) GuestFileWrite(fname string, content string) error {
	if v.vmtype == "lxc" {
		servercfg, ok := v.node.proxmox.sshcfg[v.node.id]
		if !ok {
			return errors.New("guest commands not available for LXC without SSH to hypervisor")
		}
		containerId, _ := strconv.Atoi(v.id)
		return proxmoxssh.PctPutFile(servercfg, containerId, fname, content)
	}

	var reqopt = grequests.RequestOptions{
		Data: map[string]string{
			"file":    fname,
			"content": content,
		},
	}
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/file-write", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), &reqopt)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}

	return nil
}
