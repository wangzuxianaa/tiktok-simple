package repository

import (
	"gorm.io/gorm"
)

type User struct {
	Id            int64  `gorm:"primarykey"`
	Username      string `gorm:"not null;unique;size:32"`
	Password      string `gorm:"not null;size:64"`
	FollowCount   int64  `gorm:"default:0"`
	FollowerCount int64  `gorm:"default:0"`
	IsFollow      bool   `gorm:"default:false"`
}

func (u *User) FindUserByName(username string) (*User, error) {
	var user User
	err := db.Debug().Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) FindUserById(Id int64) (*User, error) {
	var user User
	err := db.Debug().Where("Id = ?", Id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) CreateUser(username string, password string) (*User, error) {
	u.Username = username
	u.Password = password
	err := db.Create(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}
