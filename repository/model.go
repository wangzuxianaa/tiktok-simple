package repository

import "gorm.io/gorm"

// 存放数据表

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique;size:32"`
	Password string `gorm:"not null;size:64"`
}
