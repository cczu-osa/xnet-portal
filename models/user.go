package models

import "github.com/astaxie/beego/orm"

type User struct {
	Id           int
	Sid          string `orm:"unique;index"`
	PasswordHash string
	Info         *UserInfo `orm:"rel(one);null"`
}

type UserInfo struct {
	Id     int
	Sid    string `orm:"unique;index"`
	Name   string
	School string
	Major  string
}

func init() {
	orm.RegisterModel(new(User), new(UserInfo))
}
