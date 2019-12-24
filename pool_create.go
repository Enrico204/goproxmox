package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
)

func (p *proxmoxImpl) PoolCreate(name string, comment string) error {
	resp, err := p.session.Post(p.serverURL+"/api2/json/pools", &grequests.RequestOptions{
		Data: map[string]string{
			"poolid":  name,
			"comment": comment,
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
