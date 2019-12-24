package goproxmox

import (
	"errors"
)

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
