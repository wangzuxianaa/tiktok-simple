package service

import "golang.org/x/net/context"

const (
	MaxVideoNum     = 30
	UserIdNotFound  = -1
	NoGenerateToken = ""
	Add             = "1"
	Sub             = "2"
)

var ctx = context.Background()

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoMessage struct {
	Id            int64       `json:"id,omitempty"`
	Author        UserMessage `json:"author"`
	PlayUrl       string      `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount uint        `json:"favorite_count,omitempty"`
	CommentCount  uint        `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
}

type CommentMessage struct {
	Id         int64       `json:"id,omitempty"`
	User       UserMessage `json:"user"`
	Content    string      `json:"content,omitempty"`
	CreateDate string      `json:"create_date,omitempty"`
}

type UserMessage struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint   `json:"follow_count,omitempty"`
	FollowerCount uint   `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
