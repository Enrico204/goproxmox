package goproxmox

import (
	"errors"
)

func (p *proxmoxImpl) PoolDelete(name string) error {
	resp, err := p.session.Delete(p.serverURL+"/api2/json/pools/"+name, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
