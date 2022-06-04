package repository

import (
	"gorm.io/gorm"
)

//
// Comment
// @Description: 评论的数据表model，通过User和UserId使得评论和用户之间构建起约束关系，是belongs to的关系
// Deleted记录的是删除的时间，是软删除
//
type Comment struct {
	Id         int64 `gorm:"primarykey"`
	Content    string
	CreateDate string
	VideoId    int64
	UserId     int64
	User       User
	Deleted    gorm.DeletedAt
}

//
// CreateComment
// @Description: 创建一条评论信息
// @receiver c
// @return error
//
func (c *Comment) CreateComment() error {
	err := db.Debug().Create(c).Error
	return err
}

//
// DeleteComment
// @Description: 删除一条评论信息
// @receiver c
// @return error
//
func (c *Comment) DeleteComment() error {
	err := db.Debug().Delete(&Comment{}, c.Id).Error
	return err
}

//
// FindCommentsByVideoId
// @Description: 根据videoId查找评论信息，并将评论的用户信息预加载出来
// @receiver c
// @return []*Comment
// @return error
//
func (c *Comment) FindCommentsByVideoId() ([]*Comment, error) {
	var comments []*Comment
	tx := db.Debug().Preload("User").Where("video_id = ?", c.VideoId).Order("create_date desc").Find(&comments)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
