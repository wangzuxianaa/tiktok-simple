package service

import (
	"github.com/wangzuxianaa/tiktok-simple/model"
)

//
// GetCommentList
// @Description: 获取评论列表,userIdFromToken是查看评论列表人的Id
// @param videoId
// @return *[]CommentMessage
// @return error
//
func GetCommentList(videoId int64, userIdFromToken int64) (*[]CommentMessage, error) {
	// 根据视频Id查询视频的所有评论
	comments, err := model.NewCommentDaoInstance().FindCommentsByVideoId(videoId)
	if err != nil {
		return nil, err
	}

	var commentList []CommentMessage
	for _, comment := range *comments {
		// 判断是否关注了发布评论的人
		isFollow, err := model.NewFollowDaoInstance().FindFollow(comment.UserId, userIdFromToken)
		if err != nil {
			return nil, err
		}
		commentMessage := CommentMessage{
			Id: comment.Id,
			User: UserMessage{
				Id:            comment.UserId,
				Name:          comment.User.Username,
				FollowCount:   comment.User.FollowCount,
				FollowerCount: comment.User.FollowerCount,
				IsFollow:      isFollow,
			},
			Content:    comment.Content,
			CreateDate: comment.CreateDate,
		}
		commentList = append(commentList, commentMessage)
	}
	return &commentList, nil
}
