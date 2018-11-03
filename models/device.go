package models

import "github.com/astaxie/beego/orm"

type Device struct {
	Id      int
	Address string `orm:"unique;index"`
	Name    string
	User    *User `orm:"rel(fk)"`

	IPAssignments []string `orm:"-"`
}

func init() {
	orm.RegisterModel(new(Device))
}
