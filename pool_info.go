package goproxmox

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

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
