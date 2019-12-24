package goproxmox

import (
	"errors"
	"fmt"
)

type MemberStatus struct {
	VMID      string      `json:"vmid"`
	Name      string      `json:"name"`
	PID       string      `json:"pid"`
	CPUs      int         `json:"cpus"`
	CPU       float64     `json:"cpu"`
	Mem       float64     `json:"mem"`
	MaxMem    int64       `json:"maxmem"`
	Swap      float64     `json:"swap"`
	MaxSwap   int64       `json:"maxswap"`
	Uptime    int         `json:"uptime"`
	Disk      interface{} `json:"disk"`    // Sometimes it's an empty string?
	MaxDisk   interface{} `json:"maxdisk"` // Sometimes it's an empty string?
	DiskRead  int64       `json:"diskread"`
	DiskWrite int64       `json:"diskwrite"`
	Lock      string      `json:"lock"`
	Status    string      `json:"status"`
	Type      string      `json:"type"` // Sometimes it's an empty string?
	HA        struct {
		Managed int `json:"managed"`
	} `json:"ha"`
	NetIn  int64 `json:"netin"`
	NetOut int64 `json:"netout"`

	Node *string `json:"node"`

	// LXC only
	// TODO: Template  string      `json:"template"`
}

func (v *vbaseimpl) Status() (*MemberStatus, error) {
	resp, err := v.node.proxmox.session.Get(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/status/current", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var ret struct {
		Data MemberStatus `json:"data"`
	}
	err = resp.JSON(&ret)

	return &ret.Data, err
}
