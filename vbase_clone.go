package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

func (v *vbaseimpl) Clone(newhostname string, pool string, full bool, newNodeName string, timeout time.Duration) (string, error) {
	newVmId, err := v.node.proxmox.NextID()
	if err != nil {
		return "", err
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

	if newNodeName != "" {
		reqbody["target"] = newNodeName
	}

	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/clone", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		&grequests.RequestOptions{
			Data: reqbody,
		})
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.RawResponse.Status)
	}

	ret := map[string]string{}
	err = resp.JSON(&ret)
	if err != nil {
		return "", err
	}

	return newVmId, v.node.WaitForTask(ret["data"], timeout)
}
