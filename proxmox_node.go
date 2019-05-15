package goproxmox

import (
	"errors"
	"time"
)

type Node interface {
	ListLXC() ([]string, error)
	GetLXC(lxcid string) VBase
	NewLXC(lxc LXC, timeout time.Duration) (string, error)

	ListVM() ([]string, error)
	GetVM(vmid string) VBase

	ListNetworks() ([]Network, error)
	//GetNetwork(networkid string) Network
	RevertNetworkChanges(timeout time.Duration) error
	ReloadNetworkConfig(timeout time.Duration) error
	CreateNetworkConfig(network Network, timeout time.Duration) error
	UpdateNetwork(network Network, timeout time.Duration) error
	DeleteNetwork(network Network, timeout time.Duration) error

	WaitForTask(taskid string, timeout time.Duration) error
}

type nodeImpl struct {
	proxmox *proxmoxImpl `json:"-"`
	id      string       `json:"vmid"`
}

func (n *nodeImpl) GetLXC(lxcid string) VBase {
	return &vbaseimpl{vmtype: "lxc", id: lxcid, node: n}
}

func (n *nodeImpl) GetVM(vmid string) VBase {
	return &vbaseimpl{vmtype: "qemu", id: vmid, node: n}
}

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

func (n *nodeImpl) WaitForTask(taskid string, timeout time.Duration) error {
	starttime := time.Now()
	for timeout <= 0 || (time.Now().Sub(starttime) < timeout) {
		resp, err := n.proxmox.session.Get(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/tasks/"+taskid+"/status", nil)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return errors.New(resp.RawResponse.Status)
		}
		status := map[string]interface{}{}
		err = resp.JSON(&status)
		if err != nil {
			return err
		}
		if status["data"].(map[string]interface{})["status"].(string) != "running" {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return errors.New("Timeout while waiting for the operation to complete")
}
