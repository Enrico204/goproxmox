package goproxmox

import (
	"errors"
	"time"
)

func (n *nodeImpl) ReloadNetworkConfig(timeout time.Duration) error {
	resp, err := n.proxmox.session.Put(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}

	ret := map[string]string{}
	err = resp.JSON(&ret)
	if err != nil {
		return err
	}

	return n.WaitForTask(ret["data"], timeout)
}
