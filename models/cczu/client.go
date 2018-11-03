package cczu

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	CookieJar  http.CookieJar
	HttpClient *http.Client
	Data       map[string]interface{}
}

func NewClient() *Client {
	client := &Client{}
	client.CookieJar, _ = cookiejar.New(&cookiejar.Options{})
	client.HttpClient = &http.Client{Jar: client.CookieJar}
	client.Data = map[string]interface{}{}
	return client
}

func (c *Client) Login(sid, password string) (ok bool, err error) {
	ok, err = c.SsoLogin(sid, password)
	if !ok {
		return
	}

	ok, err = c.PortalLogin()
	return
}

func (c *Client) GetBasicInfo() (BasicInfo, error) {
	return c.PortalGetBasicInfo()
}
