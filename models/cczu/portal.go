package cczu

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	portalUrl           = "http://s.cczu.edu.cn/"
	ssoServicePortalUrl = "http://sso.cczu.edu.cn/sso/login?service=http%3A%2F%2Fs.cczu.edu.cn%2F"
)

// PortalLogin logs the user in s.cczu.edu.cn.
// It must be called only after a successful SsoLogin call.
func (c *Client) PortalLogin() (ok bool, err error) {
	httpClient := c.HttpClient
	res, err := httpClient.Get(ssoServicePortalUrl)
	if err != nil {
		return
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	resHtml := string(resData)

	if strings.Contains(resHtml, "欢迎登录") {
		basicInfo, err := portalParseBasicInfo(resHtml)
		if err == nil {
			c.Data["basic_info"] = basicInfo
		}

		ok = true
		err = nil
	} else {
		ok = false
		err = errors.New("failed to login")
	}
	return
}

func (c *Client) PortalGetBasicInfo() (basicInfo BasicInfo, err error) {
	if c.Data["basic_info"] != nil {
		basicInfo = c.Data["basic_info"].(BasicInfo)
		return
	}

	res, err := c.HttpClient.Get(portalUrl)
	if err != nil {
		return
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	basicInfo, err = portalParseBasicInfo(string(resData))
	if err != nil {
		return
	}

	c.Data["basic_info"] = basicInfo
	return
}

func portalParseBasicInfo(html string) (info BasicInfo, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	ok := false
	doc.Find(".person p").Each(func(i int, selection *goquery.Selection) {
		switch i {
		case 0:
			info.Name = selection.Find("i").Text()
		case 1:
			info.Sid = strings.SplitN(selection.Text(), "：", 2)[1]
		case 2:
			info.School = strings.SplitN(selection.Text(), "：", 2)[1]
		case 3:
			info.Major = strings.SplitN(selection.Text(), "：", 2)[1]
		}
		ok = i == 3
	})
	if !ok {
		err = errors.New("failed to parse student basic info")
	}
	return
}
