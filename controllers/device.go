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
		flash.Error("设备地址格式不正确，应为十位字母和数字")
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
			url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)
			body, _ := json.Marshal(map[string]interface{}{
				"authorized": true,
			})
			res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err == nil && res.StatusCode == 200 {
				succeeded = true
			}

			if succeeded {
				o.Commit()
				flash.Notice("添加成功")
			} else {
				o.Rollback()
				flash.Error("添加失败，请检查设备地址后重试")
			}
		}
	}

	flash.Store(&c.Controller)
	c.Redirect("/", 302)
}

type EditDeviceController struct {
	beego.Controller
}

func (c *EditDeviceController) Post() {
	user := MustGetSessionUser(&c.Controller)
	flash := beego.NewFlash()

	address := strings.TrimSpace(c.GetString("address"))
	name := strings.TrimSpace(c.GetString("name"))

	o := orm.NewOrm()
	device := &models.Device{Address: address, User: user}
	err := o.Read(device, "Address", "User")
	if err != nil {
		flash.Error("没有找到设备")
	} else {
		device.Name = name
		o.Update(device, "Name")
		flash.Notice("修改成功")
	}

	flash.Store(&c.Controller)
	c.Redirect("/", 302)
}

type RemoveDeviceController struct {
	beego.Controller
}

func (c *RemoveDeviceController) Post() {
	user := MustGetSessionUser(&c.Controller)
	flash := beego.NewFlash()

	address := strings.TrimSpace(c.GetString("address"))

	o := orm.NewOrm()
	device := &models.Device{Address: address, User: user}
	err := o.Read(device, "Address", "User")
	if err != nil {
		flash.Error("没有找到设备")
	} else {
		succeeded := false
		o.Begin()
		o.Delete(device)

		url := zerotierCtlApi("/network/" + zerotierNetworkId + "/member/" + address)

		// req, _ := http.NewRequest(http.MethodDelete, url, nil)
		// res, err := http.DefaultClient.Do(req)
		// if err == nil && res.StatusCode == 200 {
		// 	succeeded = true
		// }

		// Set field "authorized" to false instead of deleting it,
		// because the ZeroTier's API seems to have problem deleting nodes
		body, _ := json.Marshal(map[string]interface{}{
			"authorized": false,
		})
		res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err == nil && res.StatusCode == 200 {
			succeeded = true
		}

		if succeeded {
			o.Commit()
			flash.Notice("移除成功")
		} else {
			o.Rollback()
			flash.Error("移除失败，请稍后重试")
		}
	}

	flash.Store(&c.Controller)
	c.Redirect("/", 302)
}
