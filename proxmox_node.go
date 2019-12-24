package goproxmox

import (
	"errors"
	"time"
)

type NodeStatus struct {
	SupportLevel    string  `json:"level"`
	CPUUsagePercent float32 `json:"cpu"`
	MaxCPU          int     `json:"maxcpu"`
	MaxMem          int64   `json:"maxmem"`
	Mem             int64   `json:"mem"`
	NodeName        string  `json:"node"`
	SSLFingerprint  string  `json:"ssl_fingerprint"`
	Status          string  `json:"status"`
	UptimeSeconds   int     `json:"uptime"`
}

type Node interface {
	GetStatus() (NodeStatus, error)

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

func (n *nodeImpl) GetStatus() (NodeStatus, error) {
	nodes, err := n.proxmox.GetNodeList()
	if err != nil {
		return NodeStatus{}, err
	}
	for _, v := range nodes {
		if v.NodeName == n.id {
			return v, nil
		}
	}
	return NodeStatus{}, errors.New("can't find this node in the cluster")
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
		if resp.StatusCode == 599 {
			// Too many requests, wait more
			time.Sleep(1 * time.Second)
			continue
		} else if resp.StatusCode >= 400 {
			return errors.New(resp.RawResponse.Status)
		}
		status := map[string]interface{}{}
		err = resp.JSON(&status)
		if err != nil {
			return err
		}

		if status["data"].(map[string]interface{})["status"].(string) != "running" {
			if status["data"].(map[string]interface{})["exitstatus"].(string) != "OK" {
				return errors.New(status["data"].(map[string]interface{})["exitstatus"].(string))
			}
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("Timeout while waiting for the operation to complete")
}
