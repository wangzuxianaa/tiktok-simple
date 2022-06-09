package service

import (
	"github.com/RaymondCode/simple-demo/cache"
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/RaymondCode/simple-demo/model"
	"strconv"
	"strings"
)

//
// GetPublishList
// @Description: 获取发布视频的列表
// @param userId
// @return *[]VideoMessage
// @return error
//
func GetPublishList(userId int64) (*[]VideoMessage, error) {
	// 根据用户id查找所有发布的视频
	user, err := model.NewUserDaoInstance().FindVideosByUserId(userId)
	if err != nil {
		return nil, err
	}
	var videoList []VideoMessage
	for _, video := range user.VideoList {
		// 判断视频是否被该用户喜欢
		isLiked := cache.IsLikedByUser(userId, video.Id)
		// 获取评论总数，评论总数为redis中的评论总数和mysql中的评论总数相加
		commentCount := cache.GetCountVal(video.Id, "comment_count", int64(video.CommentCount))
		// 获取喜欢的总数，喜欢总数为redis中保存的喜欢总数和mysql中的喜欢的总数量相加
		favoriteCount := cache.GetCountVal(video.Id, "favourite_count", int64(video.FavouriteCount))
		videoMessage := VideoMessage{
			Id: video.Id,
			Author: UserMessage{
				Id:            userId,
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      user.IsFollow,
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

//
// GetFavouriteList
// @Description: 获取喜好视频的列表
// @param userId
// @return *[]VideoMessage
// @return error
//
func GetFavouriteList(userId int64) (*[]VideoMessage, error) {
	videosIdStr := cache.FindLikedVideoByUserId(userId)

	var favouriteList []VideoMessage

	for _, videoIdStr := range videosIdStr {
		var videoId int64
		var err error
		videoId, err = strconv.ParseInt(videoIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
		var video *model.Video
		// 判断视频是否被该用户喜欢
		isLiked := cache.IsLikedByUser(userId, videoId)
		// 根据视频ID查询视频信息
		video, err = model.NewVideoDaoInstance().FindVideoByVideoId(videoId)
		if err != nil {
			return nil, err
		}
		// 获取评论总数，评论总数为redis中的评论总数和mysql中的评论总数相加
		commentCount := cache.GetCountVal(videoId, "comment_count", int64(video.CommentCount))
		// 获取喜欢的总数，喜欢总数为redis中保存的喜欢总数和mysql中的喜欢的总数量相加
		favoriteCount := cache.GetCountVal(videoId, "favourite_count", int64(video.FavouriteCount))
		favourite := VideoMessage{
			Id: videoId,
			Author: UserMessage{
				Id:            video.UserId,
				Name:          video.User.Username,
				FollowCount:   video.User.FollowCount,
				FollowerCount: video.User.FollowerCount,
				IsFollow:      video.User.IsFollow,
			},
			PlayUrl:       strings.Join([]string{conf.Conf.VideoAddr, video.PlayName}, "/"),
			CoverUrl:      strings.Join([]string{conf.Conf.CoverAddr, video.CoverName}, "/"),
			FavoriteCount: uint(commentCount),
			CommentCount:  uint(favoriteCount),
			IsFavorite:    isLiked,
		}
		favouriteList = append(favouriteList, favourite)
	}
	return &favouriteList, nil
}
