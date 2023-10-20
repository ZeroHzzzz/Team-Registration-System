// 搭建数据库

package database

import (
	"backend/app/model"
	"backend/config/config"

	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	user := config.Config.GetString("database.user")
	pass := config.Config.GetString("database.pass")
	port := config.Config.GetString("database.port")
	host := config.Config.GetString("database.host")
	name := config.Config.GetString("database.name")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connect failed: ", err)
	}
	// 建表
	err1 := db.AutoMigrate(
		&model.Message{},
		&model.Usermodel{},
		&model.Team{},
		&model.Admin{},
	)

	if err1 != nil {
		log.Fatal("Database migrate failed: ", err1)
	}
	DB = db
}
