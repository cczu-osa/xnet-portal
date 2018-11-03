package main

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/cczu-osa/xnet-portal/routers"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDataBase("default", "sqlite3", "file:data/main.db")
	orm.DefaultTimeLoc = time.FixedZone("UTC+8", 8*60*60)
}

func main() {
	orm.RunSyncdb("default", false, false)
	beego.Run()
}
