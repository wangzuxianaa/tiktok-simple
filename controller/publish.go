package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
