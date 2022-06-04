package controller

import "github.com/RaymondCode/simple-demo/repository"

var video repository.Video

var DemoVideos = []VideoMessage{
	{
		Id:            5,
		Author:        DemoUser,
		PlayUrl:       "http://192.168.1.4:8080/static/6_VIDEO_20220528_075158453.mp4",
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []CommentMessage{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = UserMessage{
	Id:            7,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
