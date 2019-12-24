package goproxmox

import (
	"errors"
	"time"
)

func (p *proxmoxImpl) PoolDeleteRecursive(name string, timeout time.Duration) error {
	poolinfo, err := p.PoolInfo(name)
	if err != nil {
		return err
	}
	for _, memberinfo := range poolinfo.Members {
		var memberobj VBase
		if memberinfo.Type == "lxc" {
			memberobj = p.GetNode(*memberinfo.Node).GetLXC(memberinfo.VMID)
		} else {
			memberobj = p.GetNode(*memberinfo.Node).GetVM(memberinfo.VMID)
		}

		if memberinfo.Status == "running" {
			err = memberobj.Stop(timeout)
			if err != nil {
				return err
			}
		}
		err = memberobj.Delete(true, timeout)
		if err != nil {
			return err
		}
	}

	resp, err := p.session.Delete(p.serverURL+"/api2/json/pools/"+name, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
