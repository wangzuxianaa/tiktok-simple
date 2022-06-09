package cache

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"log"
	"strconv"
)

var Ctx = context.Background()

func GetRedisKey(videoId int64, modelName string) string {
	redisKey := fmt.Sprintf("%d:%s", videoId, modelName)
	return redisKey
}

//
// UpdateCount
// @Description: 根据actionType来更新对应键值的总数，即评论总数或者是点赞总数
// @param videoId
// @param modelName
// @param actionType
// @return int64
// @return error
//
func UpdateCount(videoId int64, modelName string, actionType string) (int64, error) {
	redisKey := GetRedisKey(videoId, modelName)
	var err error
	var Count int64

	if actionType == "1" {
		Count, err = model.RDB.Incr(Ctx, redisKey).Result()
	} else if actionType == "2" {
		Count, err = model.RDB.Decr(Ctx, redisKey).Result()
	}
	return Count, err
}

//
// GetCountVal
// @Description: 获取键值，优先获取redis中的值，若没有，获取mysql中的值
// @param Id
// @param modelName
// @param dbVal
// @return int
//
func GetCountVal(videoId int64, modelName string, dbVal int64) int64 {
	key := GetRedisKey(videoId, modelName)
	redisValStr, _ := model.RDB.Get(Ctx, key).Result()
	if redisValStr == "" {
		return dbVal
	} else {
		redisVal, err := strconv.ParseInt(redisValStr, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		return dbVal + redisVal
	}
}
