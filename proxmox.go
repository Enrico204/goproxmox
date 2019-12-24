package goproxmox

import (
	"errors"
	"time"
)

type Pool struct {
	Members []MemberStatus
	Comment string
}

type Proxmox interface {
	Login(username string, password string) error
	Logout()

	UserAdd(userid string, password string, comment string) error
	UserDelete(userid string) error

	AclAdd(path string, roles []string, users []string, groups []string, propagate bool) error
	AclDelete(path string, roles []string, users []string, groups []string) error

	PoolCreate(name string, comment string) error
	PoolDelete(name string) error
	PoolDeleteRecursive(name string, timeout time.Duration) error
	PoolList() ([]string, error)
	PoolInfo(name string) (Pool, error)

	GetNodeList() ([]NodeStatus, error)
	GetNode(nodeId string) Node

	NextID() (string, error)
}

func (p *proxmoxImpl) GetNodeList() ([]NodeStatus, error) {
	resp, err := p.session.Get(p.serverURL+"/api2/json/nodes/", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var respobj struct {
		Data []NodeStatus `json:"data"`
	}
	err = resp.JSON(&respobj)
	if err != nil {
		return nil, err
	}

	return respobj.Data, nil
}

func (p *proxmoxImpl) GetNode(nodeId string) Node {
	return &nodeImpl{id: nodeId, proxmox: p}
}
