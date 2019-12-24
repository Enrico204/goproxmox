package goproxmox

import (
	"errors"
	"fmt"
	"time"
)

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
