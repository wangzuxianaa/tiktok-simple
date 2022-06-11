package service

import "github.com/wangzuxianaa/tiktok-simple/model"

//
// GetFollowList
// @Description: 获取关注的人的列表
// @param fansId
// @param userIdFromToken
// @return *[]UserMessage
// @return error
//
func GetFollowList(fansId int64, userIdFromToken int64) (*[]UserMessage, error) {
	// 查找所有关注的人
	followIdList, err := model.NewFollowDaoInstance().AllFollow(fansId)
	if err != nil {
		return nil, err
	}
	var userMessages []UserMessage
	for _, followId := range *followIdList {
		// 查找关注人的用户信息
		followUser, err := model.NewUserDaoInstance().FindUserById(followId)
		if err != nil {
			return nil, err
		}
		// 判断操作者是否关注此人
		isFollow, err := model.NewFollowDaoInstance().FindFollow(followId, userIdFromToken)
		if err != nil {
			return nil, err
		}
		userMessage := UserMessage{
			Id:            followId,
			Name:          followUser.Username,
			FollowCount:   followUser.FollowCount,
			FollowerCount: followUser.FollowerCount,
			IsFollow:      isFollow,
		}
		userMessages = append(userMessages, userMessage)
	}
	return &userMessages, nil
}
