package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/cczu-osa/xnet-portal/models"
	"github.com/cczu-osa/xnet-portal/utils"
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
		succeeded := false
		o := orm.NewOrm()

		user := &models.User{Sid: sid}
		err := o.Read(user, "sid")

		if err != nil {
			beego.Debug("this is the first login of ", sid)

			var ok bool
			ok, user = c.loginThroughCczu(sid, password)
			if ok {
				succeeded = true
				o.Begin()
				userInfoInserted := false
				if user.Info != nil {
					_, err := o.Insert(user.Info)
					if err == nil {
						userInfoInserted = true
					}
				}
				if userInfoInserted {
					o.Insert(user)
					o.Commit()
				} else {
					o.Rollback()
				}
			}
		} else {
			beego.Debug("user ", sid, " is already in our db")
			if user.PasswordHash == utils.HashPassword(password) {
				succeeded = true
				o.Read(user.Info)
			}
		}

		if succeeded {
			beego.Debug("user ", sid, " succeeded to login, redirect now")
			c.SetSession("user", user)
			c.Redirect("/", 302)
		}
	}

	c.Get()
}

func (c *LoginController) loginThroughCczu(sid, password string) (ok bool, user *models.User) {
	client := GetSessionCczuClient(&c.Controller)
	ok, _ = client.Login(sid, password)

	if !ok {
		return
	}

	// Succeeded to login
	user = &models.User{
		Sid:          sid,
		PasswordHash: utils.HashPassword(password),
	}

	// Try to get the student's info
	basicInfo, err := client.GetBasicInfo()
	if err == nil {
		user.Info = &models.UserInfo{
			Sid:    basicInfo.Sid,
			Name:   basicInfo.Name,
			School: basicInfo.School,
			Major:  basicInfo.Major,
		}
	}

	return
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
