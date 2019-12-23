package goproxmox

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"strings"
	"time"
)

type GuestExecResult struct {
	ExitCode     int     `json:"exitcode"`
	Exited       BitBool `json:"exited"`
	OutData      string  `json:"out-data,omitempty"`
	OutTruncated BitBool `json:"out-truncated,omitempty"`
	ErrData      string  `json:"err-data,omitempty"`
	ErrTruncated BitBool `json:"err-truncated,omitempty"`
	Signal       int     `json:"signal"`
}

func (v *vbaseimpl) GuestPing() (bool, error) {
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/ping", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), nil)
	if err != nil {
		return false, err
	} else if resp.StatusCode == 500 && strings.Contains(resp.RawResponse.Status, "not running") {
		return false, nil
	} else if resp.StatusCode >= 400 {
		return false, errors.New(resp.RawResponse.Status)
	} else {
		return true, nil
	}
}

func (v *vbaseimpl) GuestExecAsync(cmd string) (uint, error) {
	var reqopt = grequests.RequestOptions{
		Data: map[string]string{
			"command": cmd,
		},
	}
	resp, err := v.node.proxmox.session.Post(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/exec", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), &reqopt)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode >= 400 {
		return 0, errors.New(resp.RawResponse.Status)
	}

	var responsejson struct {
		Data struct {
			PID uint `json:"pid"`
		} `json:"data"`
	}
	err = resp.JSON(&responsejson)
	return responsejson.Data.PID, err
}

func (v *vbaseimpl) GuestExecStatus(pid uint) (GuestExecResult, error) {
	var reqopt = grequests.RequestOptions{
		Params: map[string]string{
			"pid": fmt.Sprint(pid),
		},
	}
	resp, err := v.node.proxmox.session.Get(fmt.Sprintf("%s/api2/json/nodes/%s/%s/%s/agent/exec-status", v.node.proxmox.serverURL, v.node.id, v.vmtype, v.id), &reqopt)
	if err != nil {
		return GuestExecResult{}, err
	}
	if resp.StatusCode >= 400 {
		return GuestExecResult{}, errors.New(resp.RawResponse.Status)
	}

	var responsejson struct {
		Data GuestExecResult `json:"data"`
	}
	err = resp.JSON(&responsejson)

	return responsejson.Data, err
}

func (v *vbaseimpl) GuestExecSync(cmd string) (GuestExecResult, error) {
	pid, err := v.GuestExecAsync(cmd)
	if err != nil {
		return GuestExecResult{}, err
	}
	for {
		status, err := v.GuestExecStatus(pid)
		if err != nil {
			return GuestExecResult{}, err
		}
		if status.Exited {
			return status, nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}
