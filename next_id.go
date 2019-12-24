package goproxmox

import "errors"

func (p *proxmoxImpl) NextID() (string, error) {
	resp, err := p.session.Get(p.serverURL+"/api2/json/cluster/nextid", nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.RawResponse.Status)
	}
	status := map[string]interface{}{}
	err = resp.JSON(&status)
	if err != nil {
		return "", err
	}

	return status["data"].(string), nil
}
