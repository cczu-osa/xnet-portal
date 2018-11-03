package cczu

import (
	"errors"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ssoLoginUrl      = "http://sso.cczu.edu.cn/sso/login"
	ssoServiceUrlFmt = ssoLoginUrl + "?service=%s"
)

// SSOLogin logs the user in sso.cczu.edu.cn.
// The cookies produced by it can then be used to login other sites like
// s.cczu.edu.cn, etc.
func (c *Client) SSOLogin(sid, password string) (ok bool, err error) {
	httpClient := c.HttpClient

	res, err := httpClient.Get(ssoLoginUrl)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	formData := map[string]string{}
	doc.Find("form#fm1 input").Each(func(i int, selection *goquery.Selection) {
		name, exists := selection.Attr("name")
		if exists {
			formData[name] = selection.AttrOr("value", "")
		}
	})

	formData["username"] = sid
	formData["password"] = password

	form := url.Values{}
	for k, v := range formData {
		form.Add(k, v)
	}

	res, err = httpClient.PostForm(ssoLoginUrl, form)
	if err != nil {
		return
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if strings.Contains(string(resData), "登录成功") {
		ok = true
		err = nil
	} else {
		ok = false
		err = errors.New("failed to login")
	}
	return
}
