package model

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

//
// Video
// @Description: 视频数据表model，是one to many的关联模式，一条视频有许多评论信息。
//
type Video struct {
	Id             int64 `gorm:"primarykey"`
	UserId         int64
	User           User
	CommentList    []Comment
	PlayName       string
	CoverName      string
	FavouriteCount uint
	CommentCount   uint
	IsFavourite    bool
	Title          string
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

//
// CreateVideo
// @Description: 创建一条视频信息
// @receiver *VideoDao
// @param video
// @return error
//
func (*VideoDao) CreateVideo(video *Video) error {
	if err := DB.Debug().Create(video).Error; err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// UpdateVideoCommentCount
// @Description: 更新视频的评论总数，定期更新一次
// @receiver *VideoDao
// @param count
// @param videoId
// @return error
//
func (*VideoDao) UpdateVideoCommentCount(count int, videoId int64) error {
	err := DB.Debug().Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", count)).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// UpdateVideoFavouriteCount
// @Description: 更新视频的喜好总数定期更新一次
// @receiver *VideoDao
// @param count
// @param videoId
// @return error
//
func (*VideoDao) UpdateVideoFavouriteCount(count int, videoId int64) error {
	err := DB.Debug().Model(&Video{}).Where("id = ?", videoId).Update("favourite_count", gorm.Expr("favourite_count + ?", count)).Error
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//
// FindVideoByVideoId
// @Description: 根据视频Id查找所有的视频信息，包括视频的作者
// @receiver *VideoDao
// @param videoId
// @return *Video
// @return error
//
func (*VideoDao) FindVideoByVideoId(videoId int64) (*Video, error) {
	var video *Video
	err := DB.Debug().Preload("User").Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return video, nil
}
