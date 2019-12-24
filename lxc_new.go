package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

func (n *nodeImpl) NewLXC(lxc LXC, timeout time.Duration) (string, error) {

	if lxc.VMID == "" {
		newVmId, err := n.proxmox.NextID()
		if err != nil {
			return "", err
		}
		lxc.VMID = fmt.Sprint(newVmId)
	}

	resp, err := n.proxmox.session.Post(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/lxc", &grequests.RequestOptions{
		Data: lxc.ToMap(),
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		var errmsg struct {
			Errors map[string]string
		}
		resp.JSON(&errmsg)
		return "", errors.New(resp.RawResponse.Status + fmt.Sprint(errmsg))
	}

	ret := map[string]string{}
	err = resp.JSON(&ret)
	if err != nil {
		return "", err
	}

	return lxc.VMID, n.WaitForTask(ret["data"], timeout)
}
