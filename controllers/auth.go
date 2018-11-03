package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/cczu-osa/xnet-portal/models"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	user := c.GetSession("user")
	if user != nil {
		// The user had logged in before
		c.Redirect("/", 302)
	}
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	sid := strings.TrimSpace(c.GetString("sid"))
	password := strings.TrimSpace(c.GetString("password"))

	if len(sid) > 0 && len(password) > 0 {
		// Log into CCZU
		client := GetSessionCczuClient(&c.Controller)
		ok, _ := client.Login(sid, password)

		if ok {
			// Get the student's info
			basicInfo, err := client.PortalGetBasicInfo()
			if err == nil {
				user := models.User{
					Sid:  sid,
					Info: models.UserInfo{BasicInfo: basicInfo},
				}
				c.SetSession("user", user)
				c.Redirect("/", 302)
			}
		}
	}

	c.Get()
}

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {
	c.DelSession("user")
	DelSessionCczuClient(&c.Controller)
	c.Redirect("/", 302)
}

func (c *LogoutController) Post() {
	c.Get()
}
