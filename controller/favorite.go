package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/pkg/token"
	"github.com/RaymondCode/simple-demo/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	claims := c.MustGet("claims").(*token.Claims)
	videoIdStr := c.Query("video_id")     //视频id
	action_type := c.Query("action_type") //获取用户的点赞操作

	videoId, err := strconv.ParseInt(videoIdStr, 10, 36) //将视频Id转换为int
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	//通过videoId将视频从数据库中查找出来
	videoRepo := repository.Video{
		Id: videoId,
	}
	err = videoRepo.FindVideoByVideoId()
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//将操作者id视频播放地址连接起来作为删除标志
	deleteflag := fmt.Sprintf("%d%s", claims.UserId, videoRepo.PlayUrl)
	favorite := repository.Favorite{ //构建一个点赞列表
		Id:         videoId,       //操作的那条视频的id
		OperatorId: claims.UserId, //操作者的id
		DeleteFlag: deleteflag,
	}

	if action_type == "1" { //若为点赞操作 将点赞视频加入到点赞列表中去
		//点赞信息构建完成 将其加入点赞列表之中去
		if err := favorite.CreateFavorite(); err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" { //若是取消点赞 将视频相关信息从点赞列表之中删除
		err = favorite.DeleteFavoriteByDeleteFlag()
		if err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}

}

func FavoriteList(c *gin.Context) {
	//点赞列表应该是通过用户id从点赞列表中将此用户点赞过的列表拉取出来
	userIdStr := c.Query("user_id") //获取操作者的id
	userId, err := strconv.ParseInt(userIdStr, 10, 36)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	favorite := repository.Favorite{
		OperatorId: userId, //操作者的id
	}
	favoriteListRepos, err := favorite.FindvideosByOperationId() //通过操作者的id所有被其点赞过的视频
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//打印结果看是否正确
	//fmt.Printf("查找结果为：%#v", favoriteListRepos)

	var favorites []Video

	for _, favoriteRepo := range favoriteListRepos {
		//favoriteRepo里包含视频id---Id，操作者id---OperatorId
		//通过videoId将视频从数据库中查找出来
		videoRepo := repository.Video{
			Id: favoriteRepo.Id,
		}
		err = videoRepo.FindVideoByVideoId()
		if err != nil {
			c.JSON(http.StatusOK, CommentListResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		//videoRepo里包含了作者信息---UserId可以用于查找作者相关信息
		user := repository.User{
			Id: userId,
		}
		err = user.FindUserById()
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		favorite := Video{
			Id: favoriteRepo.Id,
			Author: User{
				Id:            videoRepo.UserId,
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
			},
			PlayUrl:       videoRepo.PlayUrl,
			CoverUrl:      videoRepo.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    true,
		}
		favorites = append(favorites, favorite)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: favorites,
	})
}
