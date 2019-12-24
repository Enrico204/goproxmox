package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

func (v *vbaseimpl) Delete(purge bool, timeout time.Duration) error {
	var reqopts = grequests.RequestOptions{}
	if purge {
		reqopts.Params = map[string]string{"purge": "1"}
	}
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
