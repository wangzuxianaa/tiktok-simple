package controller

import (
	"net/http"
	"time"

	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []VideoMessage `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var videoList []VideoMessage
	videoRepo := repository.Video{}

	videos, err := videoRepo.PullVideosFromServer()
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	for _, video := range videos {
		videoMessage := VideoMessage{
			Id: video.Id,
			Author: UserMessage{
				Id:            video.UserId,
				Name:          video.Author.Username,
				FollowCount:   video.Author.FollowCount,
				FollowerCount: video.Author.FollowerCount,
				IsFollow:      video.Author.IsFollow,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavouriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavourite,
			Title:         video.Title,
		}
		videoList = append(videoList, videoMessage)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
