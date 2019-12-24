package goproxmox

import (
	"time"
)

func (v *vbaseimpl) WaitForGuest(timeout time.Duration) (bool, error) {
	if v.vmtype == "lxc" {
		return true, nil
	}

	startts := time.Now()
	b, err := v.GuestPing()
	for err == nil && !b && time.Now().Sub(startts) < timeout {
		if !b {
			b, err = v.GuestPing()
		}
		time.Sleep(50 * time.Millisecond)
	}
	return b, err
}
