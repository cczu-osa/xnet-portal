package routers

import (
	"github.com/astaxie/beego"
	"github.com/cczu-osa/xnet-portal/controllers"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
}
