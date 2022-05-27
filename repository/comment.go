package repository

import (
	"gorm.io/gorm"
)

type Comment struct {
	Id         int64 `gorm:"primarykey"`
	Content    string
	CreateDate string
	VideoId    int64
	UserId     int64
	User       User
	Deleted    gorm.DeletedAt
}

func (c *Comment) CreateComment() error {
	err := db.Debug().Create(c).Error
	return err
}

func (c *Comment) DeleteComment() error {
	err := db.Debug().Delete(&Comment{}, c.Id).Error
	return err
}

func (c *Comment) FindCommentsByVideoId() ([]*Comment, error) {
	var comments []*Comment
	tx := db.Debug().Preload("User").Where("video_id = ?", c.VideoId).Order("create_date desc").Find(&comments)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
