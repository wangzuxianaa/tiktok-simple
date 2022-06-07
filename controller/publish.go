package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []VideoMessage `json:"video_list,omitempty"`
}

//
// Publish
// @Description: 发布视频
// @param c
//
func Publish(c *gin.Context) {
	title := c.PostForm("title")
	claims := c.MustGet("claims").(*token.Claims)
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", claims.UserId, filename)
	saveFile := filepath.Join("./public/", finalName)
	// 上传文件到指定文件夹
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 记住playurl的更新
	url := "http://192.168.43.254:8080/static"
	video := repository.Video{
		Title:   title,
		UserId:  claims.UserId,
		PlayUrl: strings.Join([]string{url, finalName}, "/"),
	}

	err = video.CreateVideo()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	videoRepo := repository.Video{
		UserId: userId,
	}
	//根据UserId查找所有的视频
	videos, err := videoRepo.FindVideosByUserId()
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	var videoList []VideoMessage

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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
