package goproxmox

import (
	"errors"
)

func (n *nodeImpl) GetStatus() (NodeStatus, error) {
	nodes, err := n.proxmox.GetNodeList()
	if err != nil {
		return NodeStatus{}, err
	}
	for _, v := range nodes {
		if v.NodeName == n.id {
			return v, nil
		}
	}
	return NodeStatus{}, errors.New("can't find this node in the cluster")
}
