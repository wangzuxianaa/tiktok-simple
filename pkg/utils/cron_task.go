package utils

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/wangzuxianaa/tiktok-simple/cache"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"log"
	"strconv"
	"strings"
)

//
// ExecuteCron
// @Description: 定时将redis中的评论总和和点赞总和数据写入mysql，每小时一次
//
func ExecuteCron() {
	c := cron.New()
	var err error
	_, err = c.AddFunc("@every 1h", func() {
		err = ScanAndUpdateCountToDB("comment_count")
		if err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		log.Print(err)
	}

	_, err = c.AddFunc("@every 1h", func() {
		err = ScanAndUpdateCountToDB("favourite_count")
		if err != nil {
			log.Print(err)
		}
	})
	if err != nil {
		log.Print(err)
	}
	c.Start()
}

//
// ScanAndUpdateCountToDB
// @Description: 扫描所有键的值，并写入mysql数据库中
// @param modelName
// @return error
//
func ScanAndUpdateCountToDB(modelName string) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = model.RDB.Scan(cache.Ctx, cursor, "*:*_count", 20).Result()
		if err != nil {
			log.Print(err)
			panic(err)
			return err
		}
		err = HandleScannedData(keys, modelName)
		if err != nil {
			log.Print(err)
			return err
		}
		if cursor == 0 {
			break
		}
	}
	log.Print("Scan and update success")
	return nil
}

//
// HandleScannedData
// @Description: 处理数据，处理完删除redis中的键值
// @param keys
// @param modelName
// @return error
//
func HandleScannedData(keys []string, modelName string) error {
	var valStr string
	var err error
	var videoId int64
	var val int
	for _, key := range keys {
		if strings.Contains(key, modelName) {
			split := strings.Split(key, ":")
			videoId, err = strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				return err
			}
			valStr, _ = model.RDB.Get(cache.Ctx, key).Result()
			if valStr == "" {
				log.Print("key does not exist")
				return errors.New("key does not exist")
			} else {
				val, err = strconv.Atoi(valStr)
				if err != nil {
					log.Print(err)
					return err
				}
			}
			switch modelName {
			// 更新到数据库
			case "comment_count":
				if err := model.NewVideoDaoInstance().UpdateVideoCommentCount(val, videoId); err != nil {
					log.Printf("updating database failed, key : %s", key)
					return err
				}
			case "favourite_count":
				if err := model.NewVideoDaoInstance().UpdateVideoFavouriteCount(val, videoId); err != nil {
					log.Printf("updating database failed, key : %s", key)
					return err
				}
			}
			// 删除缓存
			model.RDB.Del(cache.Ctx, key)
		}
	}
	log.Print("succeed to update database")
	return nil
}
