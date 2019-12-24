package goproxmox

import (
	"errors"
	"github.com/levigross/grequests"
	"net/http"
)

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
