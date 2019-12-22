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
	Comment string
}

type Proxmox interface {
	Login(username string, password string) error
	Logout()

	UserAdd(userid string, password string, comment string) error
	UserDelete(userid string) error

	AclAdd(path string, roles []string, users []string, groups []string, propagate bool) error
	AclDelete(path string, roles []string, users []string, groups []string) error

	PoolCreate(name string, comment string) error
	PoolDelete(name string) error
	PoolDeleteRecursive(name string, timeout time.Duration) error
	PoolList() ([]string, error)
	PoolInfo(name string) (Pool, error)

	GetNodeList() ([]NodeStatus, error)
	GetNode(nodeId string) Node

	NextID() (string, error)
}

func NewClient(serverURL string, verifyTLS bool, proxy string) (Proxmox, error) {
	serverURLObject, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	var greqOpts = grequests.RequestOptions{
		InsecureSkipVerify: !verifyTLS,
		RequestTimeout:     60 * time.Second,
	}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		greqOpts.Proxies = map[string]*url.URL{
			"https": proxyUrl,
		}
	}

	return &proxmoxImpl{
		serverURL:       serverURL,
		serverURLObject: serverURLObject,
		ticket:          "",
		csrf:            "",
		session:         grequests.NewSession(&greqOpts),
	}, nil
}

type proxmoxImpl struct {
	ticket          string
	csrf            string
	serverURL       string
	session         *grequests.Session
	serverURLObject *url.URL
}

func (p *proxmoxImpl) Login(username string, password string) error {

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
	p.session.HTTPClient.Jar.SetCookies(p.serverURLObject, []*http.Cookie{{
		Name:  "PVEAuthCookie",
		Value: p.ticket,
	}})
	return nil
}

func (p *proxmoxImpl) Logout() {
	p.ticket = ""
	p.csrf = ""
	p.session.RequestOptions.Headers = map[string]string{}
	p.session.HTTPClient.Jar.SetCookies(p.serverURLObject, []*http.Cookie{})
}

func (p *proxmoxImpl) GetNodeList() ([]NodeStatus, error) {
	resp, err := p.session.Get(p.serverURL+"/api2/json/nodes/", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.RawResponse.Status)
	}

	var respobj struct {
		Data []NodeStatus `json:"data"`
	}
	err = resp.JSON(&respobj)
	if err != nil {
		return nil, err
	}

	return respobj.Data, nil
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
