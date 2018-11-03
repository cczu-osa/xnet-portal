package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	beego.ReadFromRequest(&c.Controller)
	user := MustGetSessionUser(&c.Controller)
	o := orm.NewOrm()
	o.LoadRelated(user, "Devices")
	c.Data["User"] = user
	c.TplName = "index.html"
}
