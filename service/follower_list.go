package service

import "github.com/wangzuxianaa/tiktok-simple/model"

//
// GetFollowerList
// @Description: 获取粉丝列表
// @param followId
// @param userIdFromToken
// @return *[]UserMessage
// @return error
//
func GetFollowerList(followId int64, userIdFromToken int64) (*[]UserMessage, error) {
	// 获取所有粉丝
	fansIdList, err := model.NewFollowDaoInstance().AllFans(followId)
	if err != nil {
		return nil, err
	}
	var userMessages []UserMessage
	for _, fansId := range *fansIdList {
		// 查找粉丝的用户信息
		fansUser, err := model.NewUserDaoInstance().FindUserById(fansId)
		if err != nil {
			return nil, err
		}
		// 判断操作者是否关注这个粉丝
		isFollow, err := model.NewFollowDaoInstance().FindFollow(fansId, userIdFromToken)
		if err != nil {
			return nil, err
		}
		userMessage := UserMessage{
			Id:            fansId,
			Name:          fansUser.Username,
			FollowCount:   fansUser.FollowCount,
			FollowerCount: fansUser.FollowerCount,
			IsFollow:      isFollow,
		}
		userMessages = append(userMessages, userMessage)
	}
	return &userMessages, nil
}
