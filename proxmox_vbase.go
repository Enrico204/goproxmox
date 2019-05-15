package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

type VBase interface {
	Id() string

	Status() (*LXCStatus, error)
	Start(timeout time.Duration) error
	Stop(timeout time.Duration) error
	Shutdown(timeout time.Duration) error
	Delete(timeout time.Duration) error
	Clone(newhostname string, pool string, full bool, timeout time.Duration) error
	//Info() error
}

type LXCStatus struct {
	VMID      string      `json:"vmid"`
	Name      string      `json:"name"`
	Template  string      `json:"template"`
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
	Type      string      `json:"type"`
	HA        struct {
		Managed int `json:"managed"`
	} `json:"ha"`
	NetIn  int64 `json:"netin"`
	NetOut int64 `json:"netout"`

	Node *string `json:"node"`
}

type vbaseimpl struct {
	vmtype string // Can be "lxc" or "qemu"
	id     string
	node   *nodeImpl
}

func (v *vbaseimpl) Id() string {
	return v.id
}

func (v *vbaseimpl) Status() (*LXCStatus, error) {
	resp, err := v.node.proxmox.session.Get(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/status/current", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var ret struct {
		Data LXCStatus `json:"data"`
	}
	err = resp.JSON(&ret)
	return &ret.Data, err
}

func (v *vbaseimpl) Start(timeout time.Duration) error {
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/status/start", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
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

	return v.node.WaitForTask(ret["data"], timeout)
}

func (v *vbaseimpl) Stop(timeout time.Duration) error {
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/status/stop", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
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

	return v.node.WaitForTask(ret["data"], timeout)
}

func (v *vbaseimpl) Shutdown(timeout time.Duration) error {
	timeoutString := "120"
	if timeout > 0 {
		timeoutString = fmt.Sprint(int(timeout / time.Second))
	}
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/status/shutdown", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), &grequests.RequestOptions{
		Data: map[string]string{
			"timeout": timeoutString,
		},
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

	return v.node.WaitForTask(ret["data"], timeout)
}

func (v *vbaseimpl) Delete(timeout time.Duration) error {
	resp, err := v.node.proxmox.session.Delete(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
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

	return v.node.WaitForTask(ret["data"], timeout)
}

func (v *vbaseimpl) Clone(newhostname string, pool string, full bool, timeout time.Duration) error {
	newVmId, err := v.node.proxmox.NextID()
	if err != nil {
		return err
	}

	reqbody := map[string]string{
		"newid": newVmId,
	}
	if pool != "" {
		reqbody["pool"] = pool
	}

	if v.vmtype == "lxc" {
		reqbody["hostname"] = newhostname
	} else {
		reqbody["name"] = newhostname
	}

	if full {
		reqbody["full"] = "1"
	}

	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/clone", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		&grequests.RequestOptions{
			Data: reqbody,
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

	return v.node.WaitForTask(ret["data"], timeout)
}
