package goproxmox

import (
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

func (n *nodeImpl) GetLXC(lxcid string) VBase {
	return &vbaseimpl{vmtype: "lxc", id: lxcid, node: n}
}

func (n *nodeImpl) GetVM(vmid string) VBase {
	return &vbaseimpl{vmtype: "qemu", id: vmid, node: n}
}
