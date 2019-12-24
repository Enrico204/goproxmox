package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
)

func (v *vbaseimpl) DeleteNIC(id int) error {
	resp, err := v.node.proxmox.session.Put(
		fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/config", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		&grequests.RequestOptions{
			Data: map[string]string{
				"delete": fmt.Sprintf("net%d", id),
			},
		})

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
