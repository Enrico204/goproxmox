package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"time"
)

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
