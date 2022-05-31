package repository

import (
	"gorm.io/gorm"
)

type User struct { //一个用户会有多个发布的视频 也会有多个点赞的视频
	Id            int64  `gorm:"primarykey"`
	Username      string `gorm:"not null;size:32;index:idx_name_pwd"`
	Password      string `gorm:"not null;size:64;index:idx_name_pwd"`
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
	VideoList     []Video
	//LikeList []FavoriteList //外键约束
}

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

func (u *User) CreateUser() error {
	err := db.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
