package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content string
	Users   []Usermodel `gorm:"many2many:user_msg;" json:"-"`
}
