package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
	"github.com/wangzuxianaa/tiktok-simple/service"
	"net/http"
	"time"
)

//
// FeedResponse
// @Description: 视频流响应
//
type FeedResponse struct {
	service.Response
	VideoList []service.VideoMessage `json:"video_list,omitempty"`
	NextTime  int64                  `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	tokenStr, ok := c.GetQuery("token")
	var latestTime = time.Now()
	// 未登录状态
	if ok == false {
		nextTimeInt, videosList, err := service.PullVideosFromServer(0, latestTime)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: service.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response: service.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: *videosList,
			NextTime:  nextTimeInt,
		})
	} else { // 登陆状态
		claims, _ := token.ParseToken(tokenStr)
		nextTimeInt, videosList, err := service.PullVideosFromServer(claims.UserId, latestTime)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				Response: service.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response: service.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: *videosList,
			NextTime:  nextTimeInt,
		})
	}
}
