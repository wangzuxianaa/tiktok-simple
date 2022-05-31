package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/tiktok_simple?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return err
	}
	//在数据库中创建点赞表、评论表、上传视频表、用户信息表
	err = db.AutoMigrate(&Favorite{}, &Comment{}, &Video{}, &User{})
	return err
}
