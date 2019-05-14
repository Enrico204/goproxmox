package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"reflect"
	"strings"
)

type LXC struct {
	OSTemplate string `json:"ostemplate"`
	VMID       string `json:"vmid"`

	Arch               *string  `json:"arch"`
	BWLimit            *float64 `json:"bwlimit"`
	CMode              *string  `json:"cmode"`
	Console            *int     `json:"console"`
	Cores              *int     `json:"cores"`
	CPULimit           *float64 `json:"cpulimit"`
	CPUUnits           *int64   `json:"cpuunits"`
	Description        *string  `json:"description"`
	Features           *string  `json:"features"`
	Force              *int     `json:"force"`
	Hookscript         *string  `json:"hookscript"`
	Hostname           *string  `json:"hostname"`
	IgnoreUnpackErrors *int     `json:"ignore_unpack_errors" override:"ignore-unpack-errors"`
	Lock               *string  `json:"lock"`
	Memory             *int     `json:"memory"`

	MP0 *string `json:"mp0"`
	MP1 *string `json:"mp1"`
	MP2 *string `json:"mp2"`
	MP3 *string `json:"mp3"`
	MP4 *string `json:"mp4"`
	MP5 *string `json:"mp5"`
	MP6 *string `json:"mp6"`
	MP7 *string `json:"mp7"`
	MP8 *string `json:"mp8"`
	MP9 *string `json:"mp9"`

	Nameserver   *string `json:"nameserver"`
	SearchDomain *string `json:"searchdomain"`

	Net0 *string `json:"net0"`
	Net1 *string `json:"net1"`
	Net2 *string `json:"net2"`
	Net3 *string `json:"net3"`
	Net4 *string `json:"net4"`
	Net5 *string `json:"net5"`
	Net6 *string `json:"net6"`
	Net7 *string `json:"net7"`
	Net8 *string `json:"net8"`
	Net9 *string `json:"net9"`

	OnBoot     *int    `json:"onboot"`
	OsType     *string `json:"ostype"`
	Password   *string `json:"password"`
	Pool       *string `json:"pool"`
	Protection *int    `json:"protection"`
	Restore    *int    `json:"restore"`
	RootFS     *string `json:"rootfs"`

	SSH_Public_Keys *string `json:"ssh_public_keys" override:"ssh-public-keys"`
	Start           *int    `json:"start"`
	Startup         *string `json:"startup"`
	Storage         *string `json:"storage"`
	Swap            *int    `json:"swap"`
	Template        *int    `json:"template"`
	TTY             *int    `json:"tty"`
	Unique          *int    `json:"unique"`
	Unprivileged    *int    `json:"unprivileged"`
}

func (lxc *LXC) ToMap() map[string]string {
	postVars := map[string]string{}

	val := reflect.ValueOf(lxc).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		name := strings.ToLower(typeField.Name)
		if tag.Get("override") != "" {
			name = tag.Get("override")
		}
		if valueField.Kind() == reflect.String {
			postVars[name] = fmt.Sprint(valueField.Interface())
		} else if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			elem := valueField.Elem()
			postVars[name] = fmt.Sprint(elem.Interface())
		}
	}
	return postVars
}

func (n *nodeImpl) NewLXC(lxc LXC) error {
	resp, err := n.proxmox.session.Post(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/lxc", &grequests.RequestOptions{
		Data: lxc.ToMap(),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
