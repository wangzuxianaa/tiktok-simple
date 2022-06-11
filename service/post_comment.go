package service

import (
	"github.com/wangzuxianaa/tiktok-simple/cache"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"log"
	"time"
)

//
// PublishComment
// @Description: 发布评论，需要视频Id,用户Id和内容，actionType判断updateComment的方式
// @param videoId
// @param userId
// @param content
// @param actionType
// @return *model.Comment
// @return error
//
func PublishComment(videoId int64, userId int64, content string, actionType string) (*model.Comment, error) {
	comment := &model.Comment{
		VideoId:    videoId,
		UserId:     userId,
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		Content:    content,
	}
	// 更新redis中的评论总数
	if _, err := cache.UpdateCount(videoId, "comment_count", actionType); err != nil {
		return nil, err
	}

	if err := model.NewCommentDaoInstance().CreateComment(comment); err != nil {
		// 手动回滚redis中的数据
		countKey := cache.GetRedisKey(videoId, "comment_count")
		var e1 error

		if actionType == Add {
			_, e1 = model.RDB.Decr(cache.Ctx, countKey).Result()
		} else if actionType == Sub {
			_, e1 = model.RDB.Incr(cache.Ctx, countKey).Result()
		}
		if e1 != nil {
			log.Panic(e1)
		}
		return nil, err
	}

	return comment, nil
}

//
// DeleteComment
// @Description: 删除评论，需要视频Id和评论Id号，actionType判断updateComment的方式
// @param videoId
// @param commentId
// @param actionType
// @return error
//
func DeleteComment(videoId int64, commentId int64, actionType string) error {
	if err := model.NewCommentDaoInstance().DeleteComment(commentId); err != nil {
		return err
	}

	// 更新redis中的评论总数
	if _, err := cache.UpdateCount(videoId, "comment_count", actionType); err != nil {
		return err
	}
	return nil
}
