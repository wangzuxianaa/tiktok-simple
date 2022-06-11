package service

import (
	"github.com/wangzuxianaa/tiktok-simple/model"
	"log"
)

//
// FollowAction
// @Description: 关注操作
// @param followId
// @param fansId
// @return error
//
func FollowAction(followId int64, fansId int64) error {
	var err error
	var exist bool
	exist, err = model.NewFollowDaoInstance().FindFollow(followId, fansId)
	if err != nil {
		log.Print("查找失败")
		return err
	}
	if !exist {
		follow := model.Follow{
			FollowId: followId,
			FansId:   fansId,
			IsFollow: true,
		}
		if err = model.NewFollowDaoInstance().CreateFollowAndUpdateCount(&follow); err != nil {
			log.Print("创建失败")
			return err
		}
	}
	return nil
}

//
// UnfollowAction
// @Description: 取消关注操作
// @param followId
// @param fansId
// @return error
//
func UnfollowAction(followId int64, fansId int64) error {
	var err error
	var exist bool
	exist, err = model.NewFollowDaoInstance().FindFollow(followId, fansId)
	if err != nil {
		log.Print("查找失败")
		return err
	}
	if exist {
		follow := model.Follow{
			FollowId: followId,
			FansId:   fansId,
		}
		if err = model.NewFollowDaoInstance().DeleteFollowAndUpdateCount(&follow); err != nil {
			log.Print("删除失败")
			return err
		}
	}
	return nil
}
