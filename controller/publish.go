package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/model"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
	"github.com/wangzuxianaa/tiktok-simple/pkg/utils"
	"github.com/wangzuxianaa/tiktok-simple/service"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	service.Response
	VideoList []service.VideoMessage `json:"video_list"`
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
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	// 视频名
	finalVideoName := fmt.Sprintf("%d_%s", claims.UserId, filename)
	// 视频保存地址
	saveVideoFile := filepath.Join("./public/video", finalVideoName)
	// 封面名
	finalCoverName := fmt.Sprintf("%d_%s.jpeg", claims.UserId, title)
	// 封面保存地址
	savaCoverFile := filepath.Join("./public/cover", finalCoverName)
	// 上传文件到指定文件夹
	if err := c.SaveUploadedFile(data, saveVideoFile); err != nil {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 读取视频一帧作为封面
	if err := utils.ReadFrameAsJpeg(saveVideoFile, 1, savaCoverFile); err != nil {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	video := model.Video{
		Title:     title,
		UserId:    claims.UserId,
		PlayName:  finalVideoName,
		CoverName: finalCoverName,
	}

	// 创建一条视频信息
	err = model.NewVideoDaoInstance().CreateVideo(&video)
	if err != nil {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, service.Response{
		StatusCode: 0,
		StatusMsg:  finalVideoName + " uploaded successfully",
	})
}

//
// PublishList
// @Description: 视频发布列表
// @param c
//
func PublishList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	claims := c.MustGet("claims").(*token.Claims)
	var userId int64
	var err error
	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	var videoList *[]service.VideoMessage
	// 获取用户的视频发布列表
	videoList, err = service.GetPublishList(userId, claims.UserId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: *videoList,
	})
}
