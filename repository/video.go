package repository

import "gorm.io/gorm"

//
// Video
// @Description: 视频数据表model，视频表与评论表构建外键约束，是one to many的关系，一条视频有许多评论信息。
//
type Video struct {
	Id             int64 `gorm:"primarykey"`
	UserId         int64
	CommentList    []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PlayUrl        string
	CoverUrl       string
	FavouriteCount int64
	CommentCount   int64
	IsFavourite    bool
	Title          string
	Author         User
	PublishTime    string
}

//
// CreateVideo
// @Description: 创建一条视频信息
// @receiver v
// @return error
//
func (v *Video) CreateVideo() error {
	err := db.Debug().Create(v).Error
	return err
}

//
// UpdateVideoCommentCount
// @Description: 根据flag来更新评论总数，创建新的评论就加一，删除评论就减一
// @receiver v
// @param flag
// @return error
//
func (v *Video) UpdateVideoCommentCount(flag string) error {
	if flag == "1" {
		err := db.Debug().Model(v).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			return err
		}
	} else if flag == "2" {
		err := db.Debug().Model(v).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

//
// FindVideosByUserId
// @Description: 根据UserId查找视频信息
// @receiver v
// @return []*Video
// @return error
//
func (v *Video) FindVideosByUserId() ([]*Video, error) {
	var videos []*Video
	tx := db.Debug().Where("user_id = ?", v.UserId).Order("publish_time desc").Find(&videos)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (v *Video) PullVideosFromServer() ([]*Video, error) {
	var videos []*Video
	tx := db.Debug().Order("publish_time desc").Limit(30).Find(&videos)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
