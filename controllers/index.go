package controllers

import (
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	user := c.GetSession("user")
	if user == nil {
		c.Redirect("/login", 302)
	}

	c.Data["User"] = user
	c.TplName = "index.html"
}
