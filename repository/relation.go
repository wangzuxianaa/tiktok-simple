package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type Follow struct {
	FollowId int64 `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL"`
	FansId   int64 `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL"`
	IsFollow bool  `gorm:"default: false"`
}

// 查询是否有关注关系
func (u *Follow) FindFollow() (bool, error) {
	tx := db.Debug().Where("follow_id = ? AND fans_id = ?", u.FollowId, u.FansId).Find(u)
	if tx.RowsAffected != 0 {
		return true, nil
	}
	return false, nil
}

// 创建一条关注信息
func (f *Follow) CreateFollow() error {
	err := db.Create(f).Error
	if err != nil {
		return err
	}
	return nil
}

//删除一条关注信息
func (f *Follow) DeleteFollow() error {
	err := db.Where("follow_id = ? AND fans_id = ?", f.FollowId, f.FansId).Delete(f).Error
	if err != nil {
		return err
	}
	return nil
}

//所有关注的用户id
func (f *Follow) AllFollow() ([]int64, error) {
	var followIdList []int64
	var folowUsers []Follow
	tx := db.Debug().Where("fans_id = ?", f.FansId).Find(&folowUsers)
	err := tx.Error
	if err != nil {
		return followIdList, err
	}
	for _, value := range folowUsers {
		followIdList = append(followIdList, value.FollowId)
	}
	fmt.Println(followIdList)
	return followIdList, nil
}

//指定id的用户信息
func (u *User) GetFollowList(followIdList []int64) ([]User, error) {
	var followList []User
	for i := 0; i < len(followIdList); i++ {
		user := User{
			Id: followIdList[i],
		}
		err := user.FindUserById()
		if err != nil {
			return followList, err
		}
		followList = append(followList, user)
	}
	return followList, nil
}

//
// func (u *User) FindUser() error {
// 	tx := db.Debug().Where("id = ?", u.Id).Find(u)
// 	err := tx.Error
// 	if err == gorm.ErrRecordNotFound {
// 		return nil
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//所有粉丝的id
func (f *Follow) AllFans() ([]int64, error) {
	var fansIdList []int64
	var fansUsers []Follow
	tx := db.Debug().Where("follow_id = ?", f.FollowId).Find(&fansUsers)
	err := tx.Error
	if err != nil {
		return fansIdList, err
	}
	for _, value := range fansUsers {
		fansIdList = append(fansIdList, value.FansId)
	}
	fmt.Println(fansIdList)
	return fansIdList, nil
}

// update follow_count and follower_count
func (f *Follow) UpdateFollow(add bool) error {
	followUser := User{
		Id: f.FollowId,
	}
	fansUser := User{
		Id: f.FansId,
	}
	if add {
		tx := db.Model(&followUser).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
		err := tx.Error
		if err != nil {
			return err
		}
		tx = db.Model(&fansUser).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
		err = tx.Error
		if err != nil {
			return err
		}
		return nil
	} else {
		tx := db.Model(&followUser).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
		err := tx.Error
		if err != nil {
			return err
		}
		tx = db.Model(&fansUser).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
		err = tx.Error
		if err != nil {
			return err
		}
		return nil
	}
}
