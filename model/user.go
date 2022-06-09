package model

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

//
// User
// @Description: 用户数据表model
//
type User struct {
	Id            int64  `gorm:"primarykey"`
	Username      string `gorm:"not null;size:32;index:idx_name_pwd"`
	Password      string `gorm:"not null;size:64;index:idx_name_pwd"`
	FollowCount   uint
	FollowerCount uint
	IsFollow      bool
	VideoList     []Video
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

//
// FindUserByName
// @Description: 根据用户的姓名查找用户，找到了用户返回true,没有找到返回false
// @receiver *UserDao
// @param userName
// @return bool
// @return error
//
func (*UserDao) FindUserByName(userName string) (*User, bool, error) {
	var user User
	tx := DB.Debug().Select("id", "password").Where("username = ?", userName).Find(&user)
	err := tx.Error
	if err != nil {
		log.Print(err)
		return nil, false, err
	}
	if tx.RowsAffected == 0 {
		return nil, false, nil
	}
	return &user, true, nil
}

//
// FindUserById
// @Description: 根据用户的id去查找用户
// @receiver *UserDao
// @param id
// @return *User
// @return error
//
func (*UserDao) FindUserById(id int64) (*User, error) {
	var user User
	if err := DB.Debug().Where("id = ?", id).Find(&user).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

//
// CreateUser
// @Description: 创建一条用户信息
// @receiver *UserDao
// @param user
// @return error
//
func (*UserDao) CreateUser(user *User) error {
	if err := DB.Debug().Create(user).Error; err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// FindVideosByUserId
// @Description: 根据用户Id查找所有发布的视频
// @receiver *UserDao
// @param userId
// @return *User
// @return error
//
func (*UserDao) FindVideosByUserId(userId int64) (*User, error) {
	var user User
	tx := DB.Debug().Preload("VideoList", func(db *gorm.DB) *gorm.DB {
		return db.Order("videos.id desc")
	}).Where("id = ?", userId).Find(&user)
	err := tx.Error
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}
