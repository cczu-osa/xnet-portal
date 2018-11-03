package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cczu-osa/xnet-portal/models"
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

func GetSessionUser(controller *beego.Controller) *models.User {
	user := controller.GetSession("user")
	if user != nil {
		return user.(*models.User)
	}
	return nil
}

func MustGetSessionUser(controller *beego.Controller) *models.User {
	user := controller.GetSession("user")
	if user == nil {
		controller.Redirect("/login", 302)
		return nil
	}
	return user.(*models.User)
}

func SetSessionUser(controller *beego.Controller, user *models.User) {
	controller.SetSession("user", user)
}

func DelSessionUser(controller *beego.Controller) {
	controller.DelSession("user")
}
