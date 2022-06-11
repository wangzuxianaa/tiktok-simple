package cache

import (
	"fmt"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"log"
)

const (
	Add = "1"
	Sub = "2"
)

//
// GetFavoriteKey
// @Description: 根据userId生成键，用户的喜好视频id用set存储
// @param userId
// @return string
//
func GetFavoriteKey(userId int64) string {
	FavoriteKey := fmt.Sprintf("likedVideos:%d", userId)
	return FavoriteKey
}

//
// IsLikedByUser
// @Description: 判断视频是否被用户喜欢
// @param userId
// @param videoId
// @return bool
//
func IsLikedByUser(userId int64, videoId int64) bool {
	key := GetFavoriteKey(userId)
	ok, err := model.RDB.SIsMember(Ctx, key, videoId).Result()
	if err != nil {
		log.Panic(err)
	}
	return ok
}

//
// LikeActionAndUpdateCount
// @Description: 若点赞视频将视频Id加入到set中，取消点赞将视频Id从set中删除
// @param userId
// @param videoId
// @param actionType
// @return error
//
func LikeActionAndUpdateCount(userId int64, videoId int64, actionType string) (int64, error) {
	// 更新喜欢总数
	countKey := GetRedisKey(videoId, "favourite_count")
	var err error
	var Count int64

	if actionType == Add {
		Count, err = model.RDB.Incr(Ctx, countKey).Result()
	} else if actionType == Sub {
		Count, err = model.RDB.Decr(Ctx, countKey).Result()
	}
	if err != nil {
		return Count, err
	}
	// 将视频Id添加到喜欢的键值对中
	key := GetFavoriteKey(userId)
	if actionType == Add {
		_, err = model.RDB.SAdd(Ctx, key, videoId).Result()
	} else if actionType == Sub {
		_, err = model.RDB.SRem(Ctx, key, videoId).Result()
	}
	if err != nil {
		var e1 error
		// 手动回滚，redis不保证原子性
		if actionType == Add {
			Count, e1 = model.RDB.Decr(Ctx, countKey).Result()
		} else if actionType == Sub {
			Count, e1 = model.RDB.Incr(Ctx, countKey).Result()
		}
		if e1 != nil {
			return Count, e1
		}
		return Count, err
	}
	return Count, err
}

//
// FindLikedVideoByUserId
// @Description: 根据用户Id来查找喜好视频的Id
// @param userId
// @return []string
//
func FindLikedVideoByUserId(userId int64) []string {
	key := GetFavoriteKey(userId)
	likedVideos, err := model.RDB.SMembers(Ctx, key).Result()
	if err != nil {
		log.Print(err)
		return nil
	}

	return likedVideos
}
