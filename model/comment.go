package model

import (
	"log"
	"sync"
)

//
// Comment
// @Description: 评论的数据表model，通过User和UserId使得评论和用户之间构建起约束关系，是belongs to的关系
// Deleted记录的是删除的时间，是软删除
//
type Comment struct {
	Id         int64 `gorm:"primarykey"`
	Content    string
	CreateDate string `gorm:"index:idx_vid_date,priority:2"`
	VideoId    int64  `gorm:"index:idx_vid_date,priority:1"`
	UserId     int64
	User       User
}

type CommentDao struct {
}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

//
// CreateComment
// @Description: 创建一条评论信息，并更新总数
// @receiver *CommentDao
// @param comment
// @return error
//
func (*CommentDao) CreateComment(comment *Comment) error {
	err := DB.Debug().Preload("User").Create(comment).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// DeleteComment
// @Description: 删除一条评论
// @receiver *CommentDao
// @param commentId
// @return error
//
func (*CommentDao) DeleteComment(commentId int64) error {
	err := DB.Debug().Delete(&Comment{}, commentId).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// FindCommentsByVideoId
// @Description: 根据videoId查找评论信息，并将评论的用户信息预加载出来
// @receiver *CommentDao
// @param videoId
// @return []*Comment
// @return error
//
func (*CommentDao) FindCommentsByVideoId(videoId int64) (*[]*Comment, error) {
	var comments []*Comment
	tx := DB.Debug().Preload("User").Where("video_id = ?", videoId).Order("create_date desc").Find(&comments)
	err := tx.Error
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &comments, nil
}
