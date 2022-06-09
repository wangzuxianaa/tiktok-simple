package model

type Favourite struct {
	Id      int64 `gorm:"primarykey"` //视频id 唯一标识
	UserId  int64
	VideoId int64
}
