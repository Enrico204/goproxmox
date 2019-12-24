package goproxmox

import (
	"errors"
)

func (n *nodeImpl) ListLXC() ([]string, error) {
	return n.list("lxc")
}

func (n *nodeImpl) ListVM() ([]string, error) {
	return n.list("qemu")
}

func (n *nodeImpl) list(vmtype string) ([]string, error) {
	resp, err := n.proxmox.session.Get(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/"+vmtype, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var ret struct {
		Data []struct {
			VMID string `json:"vmid"`
		} `json:"data"`
	}
	err = resp.JSON(&ret)
	retstring := []string{}
	for _, lxc := range ret.Data {
		retstring = append(retstring, lxc.VMID)
	}

	return retstring, err
}
