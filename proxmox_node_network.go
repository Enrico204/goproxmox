package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"reflect"
	"strings"
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

func (net *Network) ToMap() map[string]string {
	postVars := map[string]string{}

	val := reflect.ValueOf(net).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Kind() == reflect.String {
			postVars[strings.ToLower(typeField.Name)] = fmt.Sprint(valueField.Interface())
		} else if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			elem := valueField.Elem()
			postVars[strings.ToLower(typeField.Name)] = fmt.Sprint(elem.Interface())
		}
	}
	return postVars
}

func (n *nodeImpl) CreateNetworkConfig(network Network) error {
	resp, err := n.proxmox.session.Post(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", &grequests.RequestOptions{
		Data: network.ToMap(),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}

func (n *nodeImpl) ListNetworks() ([]Network, error) {
	resp, err := n.proxmox.session.Get(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var ret struct {
		Data []Network `json:"data"`
	}
	err = resp.JSON(&ret)
	return ret.Data, err
}

func (n *nodeImpl) RevertNetworkChanges() error {
	resp, err := n.proxmox.session.Delete(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}

func (n *nodeImpl) ReloadNetworkConfig() error {
	resp, err := n.proxmox.session.Put(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}

func (n *nodeImpl) UpdateNetwork(network Network) error {
	return errors.New("Not implemented")
}

func (n *nodeImpl) DeleteNetwork(network Network) error {
	resp, err := n.proxmox.session.Delete(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/network/"+network.IFace, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
