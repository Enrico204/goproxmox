package goproxmox

import (
	"errors"
	"time"
)

type Network struct {
	Type  string `json:"type"`
	IFace string `json:"iface"`

	Active    *int      `json:"active"`
	Priority  *int      `json:"priority"`
	Families  *[]string `json:"families"`
	Autostart *int      `json:"autostart"`
	Exists    *int      `json:"exists"`
	Options   *[]string `json:"options"`
	Slaves    *string   `json:"slaves"`
	Comments  *string   `json:"comments"`
	Comments6 *string   `json:"comments6"`

	Bond_Mode             *string `json:"bond_mode"`
	Bond_Xmit_Hash_Policy *string `json:"bond_xmit_hash_policy"`

	Bridge_Ports      *string `json:"bridge_ports"`
	Bridge_Vlan_Aware *int    `json:"bridge_vlan_aware"`

	OVS_Type    *string `json:"ovs_type"`
	OVS_Ports   *string `json:"ovs_ports"`
	OVS_Bridge  *string `json:"ovs_bridge"`
	OVS_Bonds   *string `json:"ovs_bonds"`
	OVS_Options *string `json:"ovs_options"`

	Method  *string `json:"method"`
	Address *string `json:"address"`
	Netmask *string `json:"netmask"`
	Gateway *string `json:"gateway"`

	Method6  *string `json:"method6"`
	Address6 *string `json:"address6"`
	Netmask6 *string `json:"netmask6"`
	Gateway6 *string `json:"gateway6"`
}

func (n *nodeImpl) UpdateNetwork(network Network, timeout time.Duration) error {
	return errors.New("Not implemented")
}
