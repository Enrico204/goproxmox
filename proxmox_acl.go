package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
	"strings"
)

func (p *proxmoxImpl) AclAdd(path string, roles []string, users []string, groups []string, propagate bool) error {
	if len(users) == 0 && len(groups) == 0 {
		return errors.New("At least one user or group should be specified")
	}
	if len(users) != 0 && len(groups) != 0 {
		return errors.New("Only users or groups can be specified, not both")
	}

	reqbody := map[string]string{
		"path":      path,
		"roles":     strings.Join(roles, ","),
		"propagate": "0",
	}

	if len(users) > 0 {
		reqbody["users"] = strings.Join(users, ",")
	} else if len(groups) > 0 {
		reqbody["groups"] = strings.Join(groups, ",")
	}

	if propagate {
		reqbody["propagate"] = "1"
	}

	resp, err := p.session.Put(p.serverURL+"/api2/json/access/acl", &grequests.RequestOptions{
		Data: reqbody,
	})

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}

func (p *proxmoxImpl) AclDelete(path string, roles []string, users []string, groups []string) error {
	if len(users) == 0 && len(groups) == 0 {
		return errors.New("At least one user or group should be specified")
	}
	if len(users) != 0 && len(groups) != 0 {
		return errors.New("Only users or groups can be specified, not both")
	}

	reqbody := map[string]string{
		"path":   path,
		"roles":  strings.Join(roles, ","),
		"delete": "1",
	}

	if len(users) > 0 {
		reqbody["users"] = strings.Join(users, ",")
	} else if len(groups) > 0 {
		reqbody["groups"] = strings.Join(groups, ",")
	}

	resp, err := p.session.Put(p.serverURL+"/api2/json/access/acl", &grequests.RequestOptions{
		Data: reqbody,
	})

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}
	return nil
}
