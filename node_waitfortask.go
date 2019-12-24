package goproxmox

import (
	"errors"
	"time"
)

func (n *nodeImpl) WaitForTask(taskid string, timeout time.Duration) error {
	starttime := time.Now()
	for timeout <= 0 || (time.Now().Sub(starttime) < timeout) {
		resp, err := n.proxmox.session.Get(n.proxmox.serverURL+"/api2/json/nodes/"+n.id+"/tasks/"+taskid+"/status", nil)
		if err != nil {
			return err
		}
		if resp.StatusCode == 599 {
			// Too many requests, wait more
			time.Sleep(1 * time.Second)
			continue
		} else if resp.StatusCode >= 400 {
			return errors.New(resp.RawResponse.Status)
		}
		status := map[string]interface{}{}
		err = resp.JSON(&status)
		if err != nil {
			return err
		}

		if status["data"].(map[string]interface{})["status"].(string) != "running" {
			if status["data"].(map[string]interface{})["exitstatus"].(string) != "OK" {
				return errors.New(status["data"].(map[string]interface{})["exitstatus"].(string))
			}
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return errors.New("Timeout while waiting for the operation to complete")
}
