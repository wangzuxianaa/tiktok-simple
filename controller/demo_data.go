package controller

import (
	"github.com/RaymondCode/simple-demo/service"
)

var DemoVideos = []service.VideoMessage{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "http://192.168.2.10:8080/static/video/1_8224B832E4BFB4F0DBC26BFE978DA600.mp4",
		CoverUrl:      "http://192.168.2.10:8080/static/cover/1_dasdas.jpeg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []service.CommentMessage{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = service.UserMessage{
	Id:            1,
	Name:          "aaaa",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
