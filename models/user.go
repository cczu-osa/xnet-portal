package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id           int
	Sid          string `orm:"size(10);unique;index"`
	PasswordHash string
	Info         *UserInfo `orm:"rel(one);null"`
}

type UserInfo struct {
	Id     int
	Sid    string `orm:"size(10);unique;index"`
	Name   string `orm:"size(20)"`
	School string
	Major  string
}

func init() {
	orm.RegisterModel(new(User), new(UserInfo))
}
