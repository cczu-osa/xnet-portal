package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/cczu-osa/xnet-portal/models/zerotier"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	beego.ReadFromRequest(&c.Controller)
	user := MustGetSessionUser(&c.Controller)
	o := orm.NewOrm()
	o.LoadRelated(user, "Devices")

	for _, device := range user.Devices {
		if member, err := zerotier.GetMember(device.Address); err == nil {
			device.IPAssignments = member.IPAssignments
		}
	}

	c.Data["User"] = user
	c.Layout = "layout.html"
	c.TplName = "index.html"
}
