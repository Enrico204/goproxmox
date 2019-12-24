package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
	"time"
)

func (n *nodeImpl) CreateNetworkConfig(network Network, timeout time.Duration) error {
	resp, err := n.proxmox.session.Post(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", &grequests.RequestOptions{
		Data: network.ToMap(),
	})
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
