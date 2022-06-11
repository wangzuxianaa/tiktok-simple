package service

import (
	"github.com/wangzuxianaa/tiktok-simple/cache"
	"github.com/wangzuxianaa/tiktok-simple/conf"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"strconv"
	"strings"
)

//
// GetFavouriteList
// @Description: 获取用户喜好视频的列表,userId是用户Id,userIdFromToken是操作者Id
// @param userId
// @return *[]VideoMessage
// @return error
//
func GetFavouriteList(userId int64, userIdFromToken int64) (*[]VideoMessage, error) {
	// 从redis中查用户Id对应的喜好视频的Id
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
		isLiked := cache.IsLikedByUser(userIdFromToken, videoId)
		// 根据视频ID查询视频信息
		video, err = model.NewVideoDaoInstance().FindVideoByVideoId(videoId)
		if err != nil {
			return nil, err
		}
		// 获取评论总数，评论总数为redis中的评论总数和mysql中的评论总数相加
		commentCount := cache.GetCountVal(videoId, "comment_count", int64(video.CommentCount))
		// 获取喜欢的总数，喜欢总数为redis中保存的喜欢总数和mysql中的喜欢的总数量相加
		favoriteCount := cache.GetCountVal(videoId, "favourite_count", int64(video.FavouriteCount))
		// 判断操作者是否关注了视频的作者
		isFollow, err := model.NewFollowDaoInstance().FindFollow(video.UserId, userIdFromToken)
		if err != nil {
			return nil, err
		}
		favourite := VideoMessage{
			Id: videoId,
			Author: UserMessage{
				Id:            video.UserId,
				Name:          video.User.Username,
				FollowCount:   video.User.FollowCount,
				FollowerCount: video.User.FollowerCount,
				IsFollow:      isFollow,
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
