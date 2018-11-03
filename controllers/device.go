package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/cczu-osa/xnet-portal/models"
)

var (
	zerotierCtlApiRoot   = beego.AppConfig.String("zerotierctlapiroot")
	zerotierCtlAuthToken = beego.AppConfig.String("zerotierctlauthtoken")
	zerotierNetworkId    = beego.AppConfig.String("zerotiernetworkid")
)

func init() {
	zerotierCtlApiRoot = strings.TrimSuffix(zerotierCtlApiRoot, "/")
}

func zerotierCtlApi(subpath string) string {
	return zerotierCtlApiRoot + subpath + "?auth=" + zerotierCtlAuthToken
}

type AddDeviceController struct {
	beego.Controller
}

func (c *AddDeviceController) Post() {
	user := MustGetSessionUser(&c.Controller)
	flash := beego.NewFlash()

	address := strings.TrimSpace(c.GetString("address"))
	name := strings.TrimSpace(c.GetString("name"))

	if !regexp.MustCompile("[0-9a-z]{10}").MatchString(address) {
		flash.Error("设备地址格式不正确，应为 10 位字母和数字")
	} else {
		o := orm.NewOrm()

		device := &models.Device{Address: address, Name: name, User: user}
		err := o.Read(device, "Address", "User")
		if err == nil {
			// The device is already added
			flash.Error("该设备已添加过，请勿重复添加")
		} else {
			succeeded := false
			o.Begin()

			// Save device to db
			o.Insert(device)

			// Authorize the device in ZeroTier network
			body, _ := json.Marshal(map[string]interface{}{
				"authorized": true,
			})
			url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)
			res, _ := http.Post(url, "application/json", bytes.NewBuffer(body))
			if res.StatusCode == 200 {
				succeeded = true
			} else {
				flash.Error("添加失败，请检查设备地址后重试")
			}

			if succeeded {
				o.Commit()
			} else {
				o.Rollback()
			}
		}
	}

	flash.Store(&c.Controller)
	c.Redirect("/", 302)
}
