package service

import (
	"github.com/wangzuxianaa/tiktok-simple/cache"
	"github.com/wangzuxianaa/tiktok-simple/conf"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"strings"
	"time"
)

//
// FeedVideoList
// @Description: 拉取视频流，userIdFromToken
// @param userIdFromToken
// @param latestTime
// @return int64
// @return *[]VideoMessage
// @return error
//
func FeedVideoList(userIdFromToken int64, latestTime time.Time) (int64, *[]VideoMessage, error) {
	var nextTime int64
	if latestTime.IsZero() {
		latestTime = time.Now()
	}
	videos, err := model.NewVideoDaoInstance().VideoListByLimitAndTime(latestTime, MaxVideoNum)
	if err != nil {
		return time.Now().Unix(), nil, err
	}
	var videosList []VideoMessage
	for _, video := range *videos {
		// 判断视频是否被该用户喜欢
		isLiked := cache.IsLikedByUser(userIdFromToken, video.Id)
		// 获取评论总数，评论总数为redis中的评论总数和mysql中的评论总数相加
		commentCount := cache.GetCountVal(video.Id, "comment_count", int64(video.CommentCount))
		// 获取喜欢的总数，喜欢总数为redis中保存的喜欢总数和mysql中的喜欢的总数量相加
		favoriteCount := cache.GetCountVal(video.Id, "favourite_count", int64(video.FavouriteCount))
		isFollow, err := model.NewFollowDaoInstance().FindFollow(video.UserId, userIdFromToken)
		if err != nil {
			return time.Now().Unix(), nil, err
		}
		videoList := VideoMessage{
			Id: video.Id,
			Author: UserMessage{
				Id:            video.UserId,
				Name:          video.User.Username,
				FollowCount:   video.User.FollowCount,
				FollowerCount: video.User.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       strings.Join([]string{conf.Conf.VideoAddr, video.PlayName}, "/"),
			CoverUrl:      strings.Join([]string{conf.Conf.CoverAddr, video.CoverName}, "/"),
			FavoriteCount: uint(favoriteCount),
			CommentCount:  uint(commentCount),
			IsFavorite:    isLiked,
		}
		videosList = append(videosList, videoList)
	}
	if len(*videos) > 0 {
		latestTime = (*videos)[len(*videos)-1].CreatedAt
	}

	if len(*videos) == 0 {
		nextTime = time.Now().UnixNano() / 1e6
		return nextTime, &videosList, nil
	}
	nextTime = latestTime.UnixNano() / 1e6
	return nextTime, &videosList, nil
}
