package goproxmox

import "net/http"

func (p *proxmoxImpl) Logout() {
	p.ticket = ""
	p.csrf = ""
	p.session.RequestOptions.Headers = map[string]string{}
	p.session.HTTPClient.Jar.SetCookies(p.serverURLObject, []*http.Cookie{})
}
