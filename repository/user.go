package repository

import (
	"gorm.io/gorm"
)

type User struct {
	Id        int64  `gorm:"primarykey"`
	Username  string `gorm:"not null;size:32;index:idx_name_pwd"`
	Password  string `gorm:"not null;size:64;index:idx_name_pwd"`
	VideoList []Video
}

//
// FindUserByName
// @Description:根据用户的姓名查找用户，找到了用户返回true,没有找到返回false
// @receiver u
// @return bool
// @return error
//
func (u *User) FindUserByName() (bool, error) {
	tx := db.Debug().Select("id", "password").Where("username = ?", u.Username).Find(u)
	err := tx.Error
	if err != nil {
		return false, err
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

//
// FindUserById
// @Description: 根据用户的id去查找用户
// @receiver u
// @return error
//
func (u *User) FindUserById() error {
	tx := db.Debug().Where("Id = ?", u.Id).First(u)
	err := tx.Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

//
// CreateUser
// @Description: 创建一条用户信息
// @receiver u
// @return error
//
func (u *User) CreateUser() error {
	err := db.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
