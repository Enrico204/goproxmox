package goproxmox

import (
	"github.com/levigross/grequests"
	"net/url"
	"time"
)

type proxmoxImpl struct {
	ticket          string
	csrf            string
	serverURL       string
	session         *grequests.Session
	serverURLObject *url.URL
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
