package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/mitchellh/mapstructure"
	"time"
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

func (p *proxmoxImpl) PoolDeleteRecursive(name string, timeout time.Duration) error {
	poolinfo, err := p.PoolInfo(name)
	if err != nil {
		return err
	}
	for _, lxcinfo := range poolinfo.Containers {
		lxc := p.GetNode(*lxcinfo.Node).GetLXC(lxcinfo.VMID)
		if lxcinfo.Status == "running" {
			err = lxc.Stop(timeout)
			if err != nil {
				return err
			}
		}
		err = lxc.Delete(timeout)
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

func (p *proxmoxImpl) PoolList() ([]string, error) {
	resp, err := p.session.Get(p.serverURL+"/api2/json/pools", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	ret := []string{}
	var dataret struct {
		Data []struct {
			PoolId string `json:"poolid"`
		} `json:"data"`
	}
	err = resp.JSON(&dataret)
	if err != nil {
		return ret, err
	}

	for _, p := range dataret.Data {
		ret = append(ret, p.PoolId)
	}

	return ret, nil
}

func (p *proxmoxImpl) PoolInfo(name string) (Pool, error) {
	ret := Pool{
		Containers: []LXCStatus{},
	}
	resp, err := p.session.Get(p.serverURL+"/api2/json/pools/"+name, nil)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode >= 400 {
		return ret, errors.New(resp.RawResponse.Status)
	}

	var dataret struct {
		Data struct {
			Members []map[string]interface{} `json:"members"`
		} `json:"data"`
	}
	err = resp.JSON(&dataret)
	if err != nil {
		return ret, err
	}
	for _, item := range dataret.Data.Members {
		if item["type"].(string) == "lxc" {

			// Fix bad API signature
			item["vmid"] = fmt.Sprint(item["vmid"])
			item["template"] = fmt.Sprint(item["template"])

			lxc := LXCStatus{}
			err = mapstructure.Decode(item, &lxc)
			if err != nil {
				return ret, err
			}
			ret.Containers = append(ret.Containers, lxc)
		}
	}

	return ret, nil
}
