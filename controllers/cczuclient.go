package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cczu-osa/xnet-portal/models/cczu"
)

func GetSessionCczuClient(controller *beego.Controller) *cczu.Client {
	client := controller.GetSession("client")
	if client == nil {
		client = cczu.NewClient()
		controller.SetSession("client", client)
	}
	return client.(*cczu.Client)
}

func DelSessionCczuClient(controller *beego.Controller) {
	controller.DelSession("client")
}
