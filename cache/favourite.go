package cache

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"log"
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
// LikeAction
// @Description: 若点赞视频将视频Id加入到set中，取消点赞将视频Id从set中删除
// @param userId
// @param videoId
// @param actionType
// @return error
//
func LikeAction(userId int64, videoId int64, actionType string) error {
	key := GetFavoriteKey(userId)
	ok := IsLikedByUser(userId, videoId)
	if actionType == "1" && ok == true {
		_, err := model.RDB.SAdd(Ctx, key, videoId).Result()
		if err != nil {
			log.Print(err)
			return err
		}
	} else if actionType == "2" && ok == false {
		_, err := model.RDB.SRem(Ctx, key, videoId).Result()
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
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
