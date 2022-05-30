package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

//
// Init
// @Description: 数据库初始化
// @return err
//
func Init() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/tiktok_simple?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return err
	}
	err = db.AutoMigrate(&Comment{}, &Video{}, &User{})
	return err
}
