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

			var ok bool // Declare this before using it, to avoid hiding "user" of outer scope
			if ok, user = c.loginThroughCczu(sid, password); ok {
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
					succeeded = true
				} else {
					o.Rollback()
				}
			}
		} else if len(user.PasswordHash) == 0 {
			beego.Debug("there is not password hash for ", sid)
			if ok, u := c.loginThroughCczu(sid, password); ok {
				user.PasswordHash = u.PasswordHash
				o.Update(user, "PasswordHash")
				o.Read(user.Info)
				succeeded = true
			}
		} else {
			beego.Debug("user ", sid, " is already in our db")
			if utils.CompareHashAndPassword(user.PasswordHash, password) {
				o.Read(user.Info)
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
	ok, _ = client.Login(sid, password)

	if !ok {
		return
	}

	// Succeeded to login
	passwordHash := utils.GeneratePasswordHash(password)
	user = &models.User{
		Sid:          sid,
		PasswordHash: passwordHash,
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
	DelSessionUser(&c.Controller)
	DelSessionCCZUClient(&c.Controller)
	c.Redirect("/", 302)
}
