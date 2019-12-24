package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
)

func (p *proxmoxImpl) UserAdd(userid string, password string, comment string) error {
	resp, err := p.session.Post(p.serverURL+"/api2/json/access/users", &grequests.RequestOptions{
		Data: map[string]string{
			"userid":   userid,
			"password": password,
			"comment":  comment,
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

func (p *proxmoxImpl) UserDelete(userid string) error {
	resp, err := p.session.Delete(p.serverURL+"/api2/json/access/users/"+userid, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
