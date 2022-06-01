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
	// 我的mysql，根据你的进行修改
	dsn := "root:root@tcp(127.0.0.1:3306)/tiktok_simple?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return err
	}
	// 建立评论表，视频表和用户表
	err = db.AutoMigrate(&Comment{}, &Video{}, &User{})
	return err
}
