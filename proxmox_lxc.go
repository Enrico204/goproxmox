package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"reflect"
	"strings"
	"time"
)

type LXC struct {
	OSTemplate string `json:"ostemplate"`
	VMID       string `json:"vmid"`

	Arch               string  `json:"arch,omitempty"`
	BWLimit            float64 `json:"bwlimit,omitempty"`
	CMode              string  `json:"cmode,omitempty"`
	Console            int     `json:"console,omitempty"`
	Cores              uint    `json:"cores,omitempty"`
	CPULimit           float64 `json:"cpulimit,omitempty"`
	CPUUnits           int64   `json:"cpuunits,omitempty"`
	Description        string  `json:"description,omitempty"`
	Features           string  `json:"features,omitempty"`
	Force              int     `json:"force,omitempty"`
	Hookscript         string  `json:"hookscript,omitempty"`
	Hostname           string  `json:"hostname,omitempty"`
	IgnoreUnpackErrors int     `json:"ignore_unpack_errors,omitempty" override:"ignore-unpack-errors"`
	Lock               string  `json:"lock,omitempty"`
	Memory             uint    `json:"memory,omitempty"`

	MP []string `json:"mp"`

	Nameserver   string `json:"nameserver,omitempty"`
	SearchDomain string `json:"searchdomain,omitempty"`

	Net []VBaseNICSettings `json:"net"`

	OnBoot     int    `json:"onboot,omitempty"`
	OsType     string `json:"ostype,omitempty"`
	Password   string `json:"password,omitempty"`
	Pool       string `json:"pool,omitempty"`
	Protection int    `json:"protection,omitempty"`
	Restore    int    `json:"restore,omitempty"`
	RootFS     string `json:"rootfs,omitempty"`

	SSH_Public_Keys string `json:"ssh_public_keys,omitempty" override:"ssh-public-keys"`
	Start           int    `json:"start,omitempty"`
	Startup         string `json:"startup,omitempty"`
	Storage         string `json:"storage,omitempty"`
	Swap            uint   `json:"swap,omitempty"`
	Template        int    `json:"template,omitempty"`
	TTY             int    `json:"tty,omitempty"`
	Unique          int    `json:"unique,omitempty"`
	Unprivileged    int    `json:"unprivileged"`
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
		if strings.Contains(tag.Get("json"), "omitempty") && valueField.IsZero() {
			continue
		}

		if valueField.Kind() == reflect.Slice {
			for i := 0; i < valueField.Len(); i++ {
				item := valueField.Index(i)
				if item.Kind() == reflect.TypeOf(VBaseNICSettings{}).Kind() {
					v := reflect.ValueOf(item.Interface()).MethodByName("ToProxmoxString").Call([]reflect.Value{reflect.ValueOf("lxc")})
					postVars[fmt.Sprintf("%s%d", name, i)] = v[0].String()
				} else if item.Kind() == reflect.String && !item.IsZero() {
					postVars[fmt.Sprintf("%s%d", name, i)] = item.String()
				}
			}
		} else if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			elem := valueField.Elem()
			postVars[name] = fmt.Sprint(elem.Interface())
		} else {
			postVars[name] = fmt.Sprint(valueField.Interface())
		}
	}
	return postVars
}

func (n *nodeImpl) NewLXC(lxc LXC, timeout time.Duration) (string, error) {

	if lxc.VMID == "" {
		newVmId, err := n.proxmox.NextID()
		if err != nil {
			return "", err
		}
		lxc.VMID = fmt.Sprint(newVmId)
	}

	resp, err := n.proxmox.session.Post(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/lxc", &grequests.RequestOptions{
		Data: lxc.ToMap(),
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		var errmsg struct {
			Errors map[string]string
		}
		resp.JSON(&errmsg)
		return "", errors.New(resp.RawResponse.Status + fmt.Sprint(errmsg))
	}

	ret := map[string]string{}
	err = resp.JSON(&ret)
	if err != nil {
		return "", err
	}

	return lxc.VMID, n.WaitForTask(ret["data"], timeout)
}
