package models

import "github.com/cczu-osa/xnet-portal/models/cczu"

type User struct {
	Sid  string
	Info UserInfo
}

type UserInfo struct {
	cczu.BasicInfo
}
