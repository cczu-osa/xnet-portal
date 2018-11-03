package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/http"
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
		// TODO: Move ZeroTier API call to package "models"
		url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + device.Address)
		res, err := http.Get(url)
		if err == nil && res.StatusCode == 200 {
			decoder := json.NewDecoder(res.Body)
			var ztMember map[string]interface{}
			err = decoder.Decode(&ztMember)
			if err == nil {
				device.IPAssignments = make([]string, 0)
				for _, ip := range ztMember["ipAssignments"].([]interface{}) {
					device.IPAssignments = append(device.IPAssignments, ip.(string))
				}
			}
		}
	}

	c.Data["User"] = user
	c.TplName = "index.html"
}
