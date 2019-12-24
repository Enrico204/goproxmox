package goproxmox

import (
	"errors"
)

func (n *nodeImpl) ListNetworks() ([]Network, error) {
	resp, err := n.proxmox.session.Get(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var ret struct {
		Data []Network `json:"data"`
	}
	err = resp.JSON(&ret)
	return ret.Data, err
}
