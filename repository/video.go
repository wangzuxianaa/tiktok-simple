package repository

type Video struct {
	Id          int64 `gorm:"primarykey"`
	UserId      int64
	CommentList []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PlayUrl     string
	CoverUrl    string
	Title       string
}

func (v *Video) FindCommentListById(videoId int64) {
	db.Preload("CommentList ")
}
