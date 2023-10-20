package model

import (
	"crypto/md5"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// import "backend/app/model"

// User
type Usermodel struct {
	Username    string `gorm:"default:普通用户"`
	Password    string `gorm:"not null" json:"-"`
	Email       string `gorm:"primary_key"`
	Type        int
	TelePhone   int
	Sign        string `gorm:"default:哼啊啊啊啊啊啊啊啊"`
	Description string `gorm:"default:这个用户很懒，没有写简介"`
	AvataraUrl  string
	TeamID      *string `gorm:"index"`
	Team        Team    `json:"-"`
	Unread      int
	Message     []Message `gorm:"many2many:user_msg;" json:"-"`
}

func Hash(email string) [16]byte {
	return md5.Sum([]byte(strings.ToLower(strings.TrimFunc(email, func(r rune) bool {
		return r == ' '
	}))))
}

func EmailToCravatarURL(email string) string {
	return fmt.Sprintf("https://cravatar.cn/avatar/%x", Hash(email))
}

func (user *Usermodel) BeforeCreate(tx *gorm.DB) error {
	fmt.Println(user.Email)
	user.AvataraUrl = EmailToCravatarURL(user.Email)
	// 要加默认头像
	return nil
}

func (user *Usermodel) BeforeSave(tx *gorm.DB) (err error) {
	// 如果用户退出了团队
	if user.TeamID == nil {
		user.Type = 0
	}
	return
}

func (user *Usermodel) BeforeDelete(tx *gorm.DB) (err error) {
	if err := tx.Model(user).Preload("Message").First(user).Error; err != nil {
		return err
	}
	err = tx.Model(user).Association("Message").Clear()
	if err != nil {
		return err
	}
	return nil
}
