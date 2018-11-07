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
	beego.ReadFromRequest(&c.Controller)
	user := GetSessionUser(&c.Controller)
	if user != nil {
		c.Redirect("/", 302)
	}
	c.Layout = "layout.html"
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	sid := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))

	if len(sid) > 0 && len(password) > 0 {
		succeeded := false
		o := orm.NewOrm()

		user := &models.User{Sid: sid}
		if err := o.Read(user, "Sid"); err != nil {
			beego.Debug("this is the first login of ", sid)
			if ok, u := c.loginThroughCczu(sid, password); ok {
				user = u
				o.Insert(user)
				succeeded = true
			}
		} else if len(user.PasswordHash) == 0 {
			beego.Debug("there is not password hash for ", sid)
			if ok, u := c.loginThroughCczu(sid, password); ok {
				user.PasswordHash = u.PasswordHash
				o.Update(user, "PasswordHash")
				succeeded = true
			}
		} else {
			beego.Debug("user ", sid, " is already in our db")
			if utils.CompareHashAndPassword(user.PasswordHash, password) {
				succeeded = true
			}
		}

		if succeeded {
			beego.Debug("user ", sid, " succeeded to login, redirect now")
			SetSessionUser(&c.Controller, user)

			flash := beego.NewFlash()
			flash.Notice("登录成功")
			flash.Store(&c.Controller)
			c.Redirect("/", 302)
		}
	}

	flash := beego.NewFlash()
	flash.Error("登录失败，用户名密码可能不正确")
	flash.Store(&c.Controller)
	c.Redirect("/login", 302)
}

func (c *LoginController) loginThroughCczu(sid, password string) (ok bool, user *models.User) {
	client := GetSessionCCZUClient(&c.Controller)
	if ok, _ = client.Login(sid, password); !ok {
		return
	}

	// Succeeded to login
	passwordHash := utils.GeneratePasswordHash(password)
	user = &models.User{
		Sid:          sid,
		PasswordHash: passwordHash,
	}
	return
}

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {
	DelSessionUser(&c.Controller)
	DelSessionCCZUClient(&c.Controller)
	c.Redirect("/", 302)
}
