package goproxmox

import (
	"errors"
	"strconv"
	"time"
)

type Node interface {
	ListLXC() ([]string, error)
	GetLXC(lxcid string) VBase

	ListVM() ([]string, error)
	GetVM(vmid string) VBase

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
	for timeout <= 0 || (time.Now().Sub(starttime) > timeout) {
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
	return errors.New("Wait timeout")
}

func (n *nodeImpl) maxItem() (int, error) {
	ret := 1
	lxclist, err := n.ListLXC()
	if err != nil {
		return 0, err
	}

	vmlist, err := n.ListVM()
	if err != nil {
		return 0, err
	}

	for _, x := range lxclist {
		vmid, err := strconv.Atoi(x)
		if err != nil {
			continue
		}

		if ret < vmid {
			ret = vmid
		}
	}

	for _, x := range vmlist {
		vmid, err := strconv.Atoi(x)
		if err != nil {
			continue
		}

		if ret < vmid {
			ret = vmid
		}
	}

	return ret, nil
}
