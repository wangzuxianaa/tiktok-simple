package controller

import (
	"github.com/wangzuxianaa/tiktok-simple/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wangzuxianaa/tiktok-simple/pkg/token"
)

type UserListResponse struct {
	service.Response
	UserList []service.UserMessage `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	claims := c.MustGet("claims").(*token.Claims)
	fansId := claims.UserId
	followIdStr := c.Query("to_user_id")
	isFollowStr := c.Query("action_type")

	followId, err := strconv.ParseInt(followIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	var isFollow bool
	if comp := strings.Compare(isFollowStr, "1"); comp == 0 {
		isFollow = true
	} else {
		isFollow = false
	}

	if isFollow {
		if err := service.FollowAction(followId, fansId); err != nil {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		}
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 0,
			StatusMsg:  "Follow Success",
		})
	} else {
		if err := service.UnfollowAction(followId, fansId); err != nil {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		}
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 0,
			StatusMsg:  "Unfollow Success",
		})
	}
	//followUser := model.User{
	//	Id: followId,
	//}
	//fans := model.User{
	//	Id: fansId,
	//}
	//
	//err = followUser.FindUserById()
	//if err != nil {
	//	c.JSON(http.StatusOK, service.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//err = fans.FindUserById()
	//if err != nil {
	//	c.JSON(http.StatusOK, service.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//
	//follow := model.Follow{
	//	FollowId: followId,
	//	FansId:   fansId,
	//}
	//if isFollow {
	//	exist, err := follow.FindFollow()
	//	if err != nil {
	//		c.JSON(http.StatusOK, service.Response{
	//			StatusCode: 1,
	//			StatusMsg:  err.Error(),
	//		})
	//		return
	//	}
	//	if !exist {
	//		//创建新的关注列
	//		follow.IsFollow = true
	//		err = follow.CreateFollow()
	//		if err != nil {
	//			c.JSON(http.StatusOK, UserLoginResponse{
	//				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
	//			})
	//			return
	//		}
	//		err = follow.UpdateFollow(true)
	//		if err != nil {
	//			c.JSON(http.StatusOK, UserLoginResponse{
	//				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
	//			})
	//			return
	//		}
	//	}
	//	c.JSON(http.StatusOK, service.Response{
	//		StatusCode: 0,
	//		StatusMsg:  "Follow success",
	//	})
	//} else {
	//	exist, err := follow.FindFollow()
	//	if err != nil {
	//		c.JSON(http.StatusOK, service.Response{
	//			StatusCode: 1,
	//			StatusMsg:  err.Error(),
	//		})
	//		return
	//	}
	//	if exist {
	//		//删除一条关注记录
	//		err = follow.DeleteFollow()
	//		if err != nil {
	//			c.JSON(http.StatusOK, service.Response{
	//				StatusCode: 1,
	//				StatusMsg:  err.Error(),
	//			})
	//			return
	//		}
	//		err = follow.UpdateFollow(false)
	//		if err != nil {
	//			c.JSON(http.StatusOK, UserLoginResponse{
	//				Response: service.Response{StatusCode: 1, StatusMsg: err.Error()},
	//			})
	//			return
	//		}
	//	}
	//	c.JSON(http.StatusOK, service.Response{
	//		StatusCode: 0,
	//		StatusMsg:  "Unfollow success",
	//	})
	//}
}

func FollowList(c *gin.Context) {
	fansIdStr := c.Query("user_id")
	claims := c.MustGet("claims").(*token.Claims)
	var err error
	var fansId int64
	fansId, err = strconv.ParseInt(fansIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	var userList *[]service.UserMessage
	userList, err = service.GetFollowList(fansId, claims.UserId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: *userList,
	})
}

func FollowerList(c *gin.Context) {
	followIdStr := c.Query("user_id")
	claims := c.MustGet("claims").(*token.Claims)
	followId, _ := strconv.ParseInt(followIdStr, 10, 64)
	userList, err := service.GetFollowerList(followId, claims.UserId)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: *userList,
	})
}
