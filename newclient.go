package goproxmox

import (
	"github.com/levigross/grequests"
	"gitlab.com/Enrico204/goproxmox/proxmoxssh"
	"net/url"
	"time"
)

type proxmoxImpl struct {
	ticket          string
	csrf            string
	serverURL       string
	session         *grequests.Session
	serverURLObject *url.URL
	sshcfg          map[string]proxmoxssh.Config
}

// Create a new Promxox client. Parameters:
// - serverURL: the server URL (eg. https://192.0.2.1:8006 )
// - verifyTLS: whether verify or not the TLS certificate (DO NOT USE IN PRODUCTION)
// - proxy: specify a proxy server to use (empty means no proxy server)
// - sshcfg: optional map of SSH configuration for hosts (the key of the map is the server name) used in some commands that are not implemented in the API
func NewClient(serverURL string, verifyTLS bool, proxy string, sshcfg map[string]proxmoxssh.Config) (Proxmox, error) {
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

	if sshcfg == nil {
		sshcfg = map[string]proxmoxssh.Config{}
	}

	return &proxmoxImpl{
		serverURL:       serverURL,
		serverURLObject: serverURLObject,
		ticket:          "",
		csrf:            "",
		session:         grequests.NewSession(&greqOpts),
		sshcfg:          sshcfg,
	}, nil
}
