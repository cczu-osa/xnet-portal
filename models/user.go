package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id           int
	Sid          string `orm:"unique;index"`
	PasswordHash string
	Devices      []*Device `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(User))
}
