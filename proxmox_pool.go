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

//func (p *proxmoxImpl) PoolList() ([]string, error) {
//	resp,err := p.session.Get(p.serverURL + "/api2/json/pools", nil)
//	if err != nil {
//		return err
//	}
//	if resp.StatusCode >= 400 {
//		return errors.New(resp.RawResponse.Status)
//	}
//	return nil
//}

//func (p *proxmoxImpl) PoolInfo(name string) (Pool, error) {
//	resp, err := p.session.Get( p.serverURL + "/api2/json/pools/" + name, nil)
//	if err != nil {
//		return nil, err
//	}
//	if resp.StatusCode >= 400 {
//		return nil, errors.New(resp.RawResponse.Status)
//	}
//	return nil, nil
//}
