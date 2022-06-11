package service

import (
	"github.com/wangzuxianaa/tiktok-simple/cache"
	"github.com/wangzuxianaa/tiktok-simple/conf"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"strings"
)

//
// GetPublishList
// @Description: 获取发布视频的列表，userId是视频用户列表
// @param userId
// @return *[]VideoMessage
// @return error
//
func GetPublishList(userId int64, userIdFromToken int64) (*[]VideoMessage, error) {
	// 根据用户id查找所有发布的视频
	user, err := model.NewUserDaoInstance().FindVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	var videoList []VideoMessage
	for _, video := range user.VideoList {
		// 判断视频是否被操作者喜欢
		isLiked := cache.IsLikedByUser(userIdFromToken, video.Id)
		// 获取评论总数，评论总数为redis中的评论总数和mysql中的评论总数相加
		commentCount := cache.GetCountVal(video.Id, "comment_count", int64(video.CommentCount))
		// 获取喜欢的总数，喜欢总数为redis中保存的喜欢总数和mysql中的喜欢的总数量相加
		favoriteCount := cache.GetCountVal(video.Id, "favourite_count", int64(video.FavouriteCount))
		// 判断是否关注了视频的作者
		isFollow, err := model.NewFollowDaoInstance().FindFollow(userId, userIdFromToken)
		if err != nil {
			return nil, err
		}
		videoMessage := VideoMessage{
			Id: video.Id,
			Author: UserMessage{
				Id:            userId,
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       strings.Join([]string{conf.Conf.VideoAddr, video.PlayName}, "/"),
			CoverUrl:      strings.Join([]string{conf.Conf.CoverAddr, video.CoverName}, "/"),
			FavoriteCount: uint(favoriteCount),
			CommentCount:  uint(commentCount),
			IsFavorite:    isLiked,
		}
		videoList = append(videoList, videoMessage)
	}
	return &videoList, nil
}
