package model

type Admin struct {
	AdminID  int    `json:"-"`
	AdminPwd string `json:"-"`
	Key      string `json:"-"`
}
