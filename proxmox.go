package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
	"net/http"
	"net/url"
	"time"
)

type Pool struct {
	Containers []LXCStatus
	//VirtualMachines []VMStatus
}

type Proxmox interface {
	Login(username string, password string) error

	UserAdd(userid string, password string, comment string) error
	UserDelete(userid string) error

	AclAdd(path string, roles []string, users []string, groups []string, propagate bool) error
	AclDelete(path string, roles []string, users []string, groups []string) error

	PoolCreate(name string, comment string) error
	PoolDelete(name string) error
	PoolDeleteRecursive(name string, timeout time.Duration) error
	PoolList() ([]string, error)
	PoolInfo(name string) (Pool, error)

	GetNode(nodeId string) Node

	NextID() (string, error)
}

func NewClient(serverURL string, verifyTLS bool) Proxmox {
	return &proxmoxImpl{
		serverURL: serverURL,
		ticket:    "",
		csrf:      "",
		session: grequests.NewSession(&grequests.RequestOptions{
			InsecureSkipVerify: !verifyTLS,
			RequestTimeout:     5 * time.Second,
		}),
	}
}

type proxmoxImpl struct {
	ticket    string
	csrf      string
	serverURL string
	session   *grequests.Session
}

func (p *proxmoxImpl) Login(username string, password string) error {
	serverURLObject, err := url.Parse(p.serverURL)
	if err != nil {
		// TODO: move to New()
		return err
	}

	resp, err := p.session.Post(p.serverURL+"/api2/json/access/ticket", &grequests.RequestOptions{
		Data: map[string]string{
			"username": username,
			"password": password,
		},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.RawResponse.Status)
	}

	loginresponse := struct {
		Data struct {
			CSRFPreventionToken string                 `json:"CSRFPreventionToken"`
			Ticket              string                 `json:"ticket"`
			Cap                 map[string]interface{} `json:"cap"`
			UserName            string                 `json:"username"`
		} `json:"data"`
	}{}
	err = resp.JSON(&loginresponse)
	if err != nil {
		return err
	}

	p.ticket = loginresponse.Data.Ticket
	p.csrf = loginresponse.Data.CSRFPreventionToken
	p.session.RequestOptions.Headers = map[string]string{
		"CSRFPreventionToken": p.csrf,
	}
	p.session.HTTPClient.Jar.SetCookies(serverURLObject, []*http.Cookie{{
		Name:  "PVEAuthCookie",
		Value: p.ticket,
	}})
	return nil
}

func (p *proxmoxImpl) GetNode(nodeId string) Node {
	return &nodeImpl{id: nodeId, proxmox: p}
}

func (p *proxmoxImpl) NextID() (string, error) {
	resp, err := p.session.Get(p.serverURL+"/api2/json/cluster/nextid", nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.RawResponse.Status)
	}
	status := map[string]interface{}{}
	err = resp.JSON(&status)
	if err != nil {
		return "", err
	}

	return status["data"].(string), nil
}
