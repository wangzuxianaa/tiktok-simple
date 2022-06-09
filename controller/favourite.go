package controller

import (
	"github.com/RaymondCode/simple-demo/cache"
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FavouriteAction(c *gin.Context) {
	claims := c.MustGet("claims").(*token.Claims)
	videoIdStr := c.Query("video_id")    //视频id
	actionType := c.Query("action_type") //获取用户的点赞操作

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36) //将视频Id转换为int
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}

	// 点赞操作
	if actionType == "1" {
		// 点赞数据存入redis
		if err := cache.LikeAction(claims.UserId, videoId, actionType); err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		// 更新redis中的点赞总数
		_, err = cache.UpdateCount(videoId, "favourite_count", actionType)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{StatusCode: 0, StatusMsg: "success"},
		})
	} else if actionType == "2" {
		// 删除点赞数据
		if err := cache.LikeAction(claims.UserId, videoId, actionType); err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		// 点赞总数减一
		_, err = cache.UpdateCount(videoId, "favourite_count", actionType)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{StatusCode: 0, StatusMsg: "success"},
		})
	}

}

func FavouriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	var userId int64
	var err error
	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	var favouriteList *[]service.VideoMessage
	// 根据用户id获取点赞列表
	favouriteList, err = service.GetFavouriteList(userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: *favouriteList,
	})
}
