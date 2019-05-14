package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
	"net/http"
	"time"
)

type Pool interface {
}

type Proxmox interface {
	Login(username string, password string) error

	UserAdd(userid string, password string, comment string) error
	UserDelete(userid string) error

	AclAdd(path string, roles []string, users []string, groups []string, propagate bool) error
	AclDelete(path string, roles []string, users []string, groups []string) error

	PoolCreate(name string, comment string) error
	PoolDelete(name string) error
	//PoolList() ([]string, error)
	//PoolInfo(name string) (Pool, error)

	GetNode(nodeId string) Node
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
		Data map[string]string `json:"data"`
	}{}
	err = resp.JSON(&loginresponse)
	if err != nil {
		return err
	}

	p.ticket = loginresponse.Data["ticket"]
	p.csrf = loginresponse.Data["CSRFPreventionToken"]
	p.session.RequestOptions.Headers["CSRFPreventionToken"] = p.csrf
	p.session.RequestOptions.Cookies = append(p.session.RequestOptions.Cookies, &http.Cookie{Name: "PVEAuthCookie", Value: p.ticket})
	return nil
}

func (p *proxmoxImpl) GetNode(nodeId string) Node {
	return &nodeImpl{id: nodeId, proxmox: p}
}
