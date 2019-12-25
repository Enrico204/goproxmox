package goproxmox

import (
	"errors"
	"fmt"
)

func (v *vbaseimpl) GetNIC(idx int) (VBaseNICSettings, error) {
	var settings = VBaseNICSettings{}
	resp, err := v.node.proxmox.session.Get(
		fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/config", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id),
		nil)

	if err != nil {
		return settings, err
	}
	if resp.StatusCode >= 400 {
		return settings, errors.New(resp.RawResponse.Status)
	}

	var respjson struct {
		Data map[string]interface{}
	}
	err = resp.JSON(&respjson)
	if err != nil {
		return settings, err
	}
	cfg, ok := respjson.Data[fmt.Sprintf("net%d", idx)]
	if !ok {
		return settings, errors.New("invalid NIC index")
	}
	settings.Id = idx

	err = settings.FromProxmoxString(cfg.(string))
	return settings, err
}
