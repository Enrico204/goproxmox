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
		Members: []MemberStatus{},
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
			Comment string                   `json:"comment"`
		} `json:"data"`
	}

	err = resp.JSON(&dataret)
	if err != nil {
		return ret, err
	}
	for _, item := range dataret.Data.Members {
		// Fix bad API signature
		item["vmid"] = fmt.Sprint(item["vmid"])
		/*if item["type"].(string) == "lxc" {
			item["template"] = fmt.Sprint(item["template"])
		}*/

		mb := MemberStatus{}
		err = mapstructure.Decode(item, &mb)
		if err != nil {
			return ret, err
		}
		ret.Members = append(ret.Members, mb)
	}
	ret.Comment = dataret.Data.Comment

	return ret, nil
}
