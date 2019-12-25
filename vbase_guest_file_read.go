package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"gitlab.com/Enrico204/goproxmox/proxmoxssh"
	"strconv"
)

func (v *vbaseimpl) GuestFileRead(fname string) (string, error) {
	if v.vmtype == "lxc" {
		servercfg, ok := v.node.proxmox.sshcfg[v.node.id]
		if !ok {
			return "", errors.New("guest commands not available for LXC without SSH to hypervisor")
		}
		containerId, _ := strconv.Atoi(v.id)
		return proxmoxssh.PctGetFile(servercfg, containerId, fname)
	}

	var reqopt = grequests.RequestOptions{
		Params: map[string]string{
			"file": fname,
		},
	}
	resp, err := v.node.proxmox.session.Get(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/file-read", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), &reqopt)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.RawResponse.Status)
	}

	var respjson struct {
		Data struct {
			Content   string `json:"content"`
			BytesRead int64  `json:"bytes-read"`
		} `json:"data"`
	}

	err = resp.JSON(&respjson)
	if err != nil {
		return "", err
	}
	return respjson.Data.Content, nil
}
