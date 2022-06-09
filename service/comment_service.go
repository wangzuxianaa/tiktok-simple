package service

import (
	"github.com/RaymondCode/simple-demo/cache"
	"github.com/RaymondCode/simple-demo/model"
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
	if err := model.NewCommentDaoInstance().CreateComment(comment); err != nil {
		return nil, err
	}

	// 更新redis中的评论总数
	if _, err := cache.UpdateCount(videoId, "comment_count", actionType); err != nil {
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

//
// GetCommentList
// @Description: 获取评论列表
// @param videoId
// @return *[]CommentMessage
// @return error
//
func GetCommentList(videoId int64) (*[]CommentMessage, error) {
	// 根据视频Id查询视频的所有评论
	comments, err := model.NewCommentDaoInstance().FindCommentsByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	var commentList []CommentMessage
	for _, comment := range comments {
		commentMessage := CommentMessage{
			Id: comment.Id,
			User: UserMessage{
				Id:            comment.UserId,
				Name:          comment.User.Username,
				FollowCount:   comment.User.FollowCount,
				FollowerCount: comment.User.FollowerCount,
				IsFollow:      comment.User.IsFollow,
			},
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		}
		commentList = append(commentList, commentMessage)
	}
	return &commentList, nil
}
