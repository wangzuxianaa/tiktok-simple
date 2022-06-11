package model

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
)

type Follow struct {
	FollowId int64 `gorm:"primaryKey;autoIncrement:false" sql:"type:INT(10) UNSIGNED NOT NULL"`
	FansId   int64 `gorm:"primaryKey;autoIncrement:false" sql:"type:INT(10) UNSIGNED NOT NULL"`
	IsFollow bool  `gorm:"default: false"`
}

type FollowDao struct {
}

var followDao *FollowDao
var followOnce sync.Once

func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

//
// FindFollow
// @Description: 查询表中是否有followId和fansId判断是否有关注关系
// @receiver *FollowDao
// @param followId
// @param fansId
// @return bool
// @return error
//
func (*FollowDao) FindFollow(followId int64, fansId int64) (bool, error) {
	var f Follow
	tx := DB.Debug().Where("follow_id = ? AND fans_id = ?", followId, fansId).Find(&f)
	err := tx.Error
	if err != nil {
		log.Panicln(err)
	}
	if tx.RowsAffected != 0 {
		return true, nil
	}
	return false, nil
}

//
// CreateFollowAndUpdateCount
// @Description: 创建记录并更新总数，用事务保证，错误回滚
// @receiver *FollowDao
// @param f
// @return error
//
func (*FollowDao) CreateFollowAndUpdateCount(f *Follow) error {
	followUser := User{
		Id: f.FollowId,
	}
	fansUser := User{
		Id: f.FansId,
	}
	// 开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(f).Error; err != nil {
			return err
		}
		if err := tx.Model(&followUser).
			UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&fansUser).
			UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}

//
// DeleteFollowAndUpdateCount
// @Description: 删除记录并更新总数，用事务保证，错误回滚
// @receiver *FollowDao
// @param f
// @return error
//
func (*FollowDao) DeleteFollowAndUpdateCount(f *Follow) error {
	followUser := User{
		Id: f.FollowId,
	}
	fansUser := User{
		Id: f.FansId,
	}
	// 开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Delete(f).Error; err != nil {
			return err
		}
		if err := tx.Model(&followUser).
			UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&fansUser).
			UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}

//
// AllFollow
// @Description: 根据某个用户关注的人
// @receiver *FollowDao
// @param fansId
// @return []int64
// @return error
//
func (*FollowDao) AllFollow(fansId int64) (*[]int64, error) {
	var followIdList []int64
	var followUsers []Follow
	tx := DB.Debug().Where("fans_id = ?", fansId).Find(&followUsers)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	for _, value := range followUsers {
		followIdList = append(followIdList, value.FollowId)
	}
	return &followIdList, nil
}

//
// AllFans
// @Description: 查找某个用户的所有粉丝列表
// @receiver *FollowDao
// @param followId
// @return []int64
// @return error
//
func (*FollowDao) AllFans(followId int64) (*[]int64, error) {
	var fansIdList []int64
	var fansUsers []Follow
	tx := DB.Debug().Where("follow_id = ?", followId).Find(&fansUsers)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	for _, value := range fansUsers {
		fansIdList = append(fansIdList, value.FansId)
	}
	fmt.Println(fansIdList)
	return &fansIdList, nil
}
